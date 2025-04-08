package main

import (
	"analytics/connectors"
	"analytics/controller"
	mongodb "analytics/repository"
	"log"
	"os"

	"github.com/labstack/echo/v4"
)

type container struct {
	analyticsLogsInstance controller.AnalyticsLogs
}

func loadContainer() *container {
	mongoCon := mongodb.MongoConnect()
	if mongoCon == nil {
		log.Fatal("Failed to connect to MongoDB")
	}
	return &container{
		analyticsLogsInstance: controller.AnalyticsLogs{MongoCon: mongodb.MongoConnect()},
	}
}

func init() {
	connectors.LoadEnv()
	connectors.LoadLogger()
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8001"
	}

	containerInstance := loadContainer()
	log.Println("Database Connected Successfully")
	e := echo.New()
	e.POST("/upload-csv", containerInstance.analyticsLogsInstance.UploadCSV)
	e.GET("/get-revenue", containerInstance.analyticsLogsInstance.GetRevenue)
	e.GET("/get-customer-analysis", containerInstance.analyticsLogsInstance.GetCustomerAndOrderStats)
	e.GET("/get-other-calculations", containerInstance.analyticsLogsInstance.GetProfitMarginByProduct)
	e.GET("/refresh", containerInstance.analyticsLogsInstance.RefreshData)
	log.Printf("Starting server on port %s", port)
	e.Logger.Fatal(e.Start(":" + port))
}
