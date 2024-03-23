package nmiCsv

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"time"
)

type fileError struct {
	err string
}

func (e *fileError) Error() string {
	return fmt.Sprintf(e.err)
}

type MeterStruct struct {
	Nmi         string
	Timestamp   time.Time
	Consumption float64
}

type MeterReadingsData struct {
	Timestamp   time.Time
	Consumption float64
}

// parse each record in csv to get the nmi and interval
func NmiParser(reader *csv.Reader) ([]MeterStruct, error) {
	// full data in the file. this will be the output to write to db
	summary := make(map[string]map[string]MeterReadingsData)
	data := make([]MeterStruct, 0)

	// get current nmi in the set
	currentNmi := ""
	// set the current interval based on each 200 record
	currentInterval := 0

	currentRow := 0

	for {
		record, err := reader.Read()

		recordNum := record[0]

		// If file does not start with 100, return error
		if currentRow == 0 && recordNum != "100" {
			return data, &fileError{
				err: "File does not start with 100",
			}
		}

		// File ends with 900
		if recordNum == "900" {
			currentRow += 1
			break
		}

		// If end of the file is encountered
		if err == io.EOF {
			return data, &fileError{
				err: "File does not end with 900",
			}
		}

		// parse 200 record
		if recordNum == "200" {
			// string to int
			interval, err := strconv.Atoi(record[8])
			if err != nil {
				// check if interval has an error converting to number
				return data, &fileError{
					err: "Convert interval error",
				}
			}
			nmi := record[1]

			// set data for parsing 300 records
			currentNmi = nmi
			currentInterval = interval

			// if nmi is not in the set, add it
			if _, ok := summary[nmi]; !ok {
				summary[nmi] = make(map[string]MeterReadingsData)
			}
		}

		// parse 300 record
		if recordNum == "300" {
			// If 200 is not before 300 record, error
			if currentNmi == "" && currentInterval == 0 {
				return data, &fileError{
					err: "No 200 record before 300",
				}
			}

			// convert first record to date (start of day)
			startDate, err := time.Parse("20060102 15:04", record[1]+" 00:00")
			if err != nil {
				// if error converting to date
				return data, &fileError{
					err: "Wrong date format in 300 record",
				}
			}

			endDate := startDate.Add(time.Duration(24) * time.Hour)
			currentIndex := 2
			for startDate.Before(endDate) {
				// convert current usage to float for addition
				usage, err := strconv.ParseFloat(record[currentIndex], 3)
				if err != nil {
					// if error converting to float
					return data, &fileError{
						err: "Consumption data is not a number",
					}
				}

				if _, ok := summary[currentNmi][startDate.Format("20060102 1504")]; !ok {
					// add usage to the summary
					summary[currentNmi][startDate.Format("20060102 1504")] = MeterReadingsData{
						Consumption: usage,
						Timestamp:   startDate,
					}
				} else {
					// if the date is already in the summary, add the usage
					summary[currentNmi][startDate.Format("20060102 1504")] = MeterReadingsData{
						Consumption: summary[currentNmi][startDate.Format("20060102 1504")].Consumption + usage,
						Timestamp:   startDate,
					}
				}

				// increment the date using interval and parse the data in the next index
				startDate = startDate.Add(time.Duration(currentInterval) * time.Minute)
				currentIndex += 1
			}
		}

		currentRow += 1
	}

	// convert summary to MeterStruct. for saving to db
	for nmi, consumptionData := range summary {
		for _, consumption := range consumptionData {
			data = append(data, MeterStruct{
				Nmi:         nmi,
				Timestamp:   consumption.Timestamp,
				Consumption: consumption.Consumption,
			})
		}
	}

	return data, nil
}
