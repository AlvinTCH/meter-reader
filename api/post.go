package api

import (
	"encoding/csv"
	"net/http"

	"log"

	"flo/database"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	nmiCsv "flo/csv"
)

// PostCSV receives a CSV file and parse it
func PostCSV(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		// log.Fatal(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
	}

	// open file
	fileData, err := file.Open()
	if err != nil {
		// log.Fatal("Error while opening the file", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
	}

	reader := csv.NewReader(fileData)
	// parse csv file
	nmiData, err := nmiCsv.NmiParser(reader)
	if err != nil {
		// log.Fatal("Error while reading the file", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
	}

	log.Println(nmiData)

	db := database.GlobDB
	dbSession := db.Session(&gorm.Session{CreateBatchSize: 1000})

	for _, data := range nmiData {
		dbSession.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "timestamp"}, {Name: "nmi"}},
			DoUpdates: clause.AssignmentColumns([]string{"consumption"}),
		}).Create(&database.MeterReadings{
			ID:          uuid.New().String(),
			Nmi:         data.Nmi,
			Timestamp:   data.Timestamp,
			Consumption: data.Consumption,
		})
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "File saved successfully!"})
}
