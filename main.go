package main

import (
	"github.com/gin-gonic/gin"

	api "flo/api"
	database "flo/database"

	_ "github.com/lib/pq"
)

func main() {
	database.Open()
	database.Migrate()

	defer database.GlobSqlDB.Close()

	router := gin.Default()

	router.POST("/upload-csv", api.PostCSV)
	router.GET("/get-all-readings", api.GetAllReadings)

	router.Run()
}
