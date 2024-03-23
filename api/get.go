package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	database "flo/database"

	"github.com/gin-gonic/gin"
)

// get specific reading
func GetAllReadings(c *gin.Context) {
	db := database.GlobDB

	// convert first record to date (start of day)
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
		log.Fatal("Error reading rows", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
	}

	for rows.Next() {
		db.ScanRows(rows, &readings)
		fmt.Println(readings)
	}

	c.IndentedJSON(http.StatusOK, result)
}
