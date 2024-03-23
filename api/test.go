package api

import (
	"net/http"
	"time"

	database "flo/database"

	"github.com/gin-gonic/gin"
)

type MeterReadingsReturn struct {
	ID          string  `json:"id"`
	Timestamp   string  `json:"timestamp"`
	Nmi         string  `json:"nmi"`
	Consumption float64 `json:"reading"`
}

// check for updating in database
func TestUpdateReadings(c *gin.Context) {
	db := database.GlobDB

	// convert record to date
	queryDate, err := time.Parse("20060102 15:04", "20050301 06:30")
	if err != nil {
		// if error converting to date
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
	}

	var readings = database.MeterReadings{}
	result := db.Where("timestamp = ?", queryDate).Find(&readings)

	rows, err := result.Rows()
	if err != nil {
		// log.Fatal("Error reading rows", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
	}

	//
	data := make([]MeterReadingsReturn, 0)
	for rows.Next() {
		db.ScanRows(rows, &readings)
		data = append(data, MeterReadingsReturn{
			ID:          readings.ID,
			Timestamp:   readings.Timestamp.Format("2006-01-02 15:04"),
			Nmi:         readings.Nmi,
			Consumption: readings.Consumption,
		})
	}

	c.IndentedJSON(http.StatusOK, data)
}
