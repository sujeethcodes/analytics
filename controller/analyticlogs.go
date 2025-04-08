package controller

import (
	"analytics/entity"
	"analytics/usecase"
	"context"
	"net/http"
	"time"

	mongodb "analytics/repository"

	"github.com/labstack/echo/v4"
)

type AnalyticsLogs struct {
	MongoCon *mongodb.MongoCon
}

func (a *AnalyticsLogs) UploadCSV(c echo.Context) error {
	use := usecase.AnalyticsLogs{MongoCon: a.MongoCon}
	ctx := context.Background()

	// Get the uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to read file"})
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to open file"})
	}
	defer src.Close()

	// Call the usecase function to process the CSV
	err = use.ProcessUploadedCSV(ctx, src)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to process CSV"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "CSV uploaded and processed successfully"})
}

func (a *AnalyticsLogs) RefreshData(c echo.Context) error {
	use := usecase.AnalyticsLogs{MongoCon: a.MongoCon}
	ctx := context.Background()

	err := use.RefreshAnalyticsData(ctx, "path/to/your/data.csv")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Data refresh failed"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Data refreshed successfully"})
}

func (a *AnalyticsLogs) GetRevenue(c echo.Context) error {
	use := usecase.AnalyticsLogs{MongoCon: a.MongoCon}

	startDateStr := c.QueryParam("start_date")
	endDateStr := c.QueryParam("end_date")
	// type is empty total revenue is return based on date range
	aggType := c.QueryParam("type") // "", "product", "category", or "region"

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid start_date format"})
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid end_date format"})
	}

	ctx := context.Background()
	result, err := use.CalculateRevenueByType(ctx, startDate, endDate, aggType)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to calculate revenue"})
	}

	return c.JSON(http.StatusOK, entity.Response{
		Status:  200,
		Message: "Data fetching Successfully",
		Data:    result,
	})
}

func (a *AnalyticsLogs) GetCustomerAndOrderStats(c echo.Context) error {
	use := usecase.AnalyticsLogs{MongoCon: a.MongoCon}

	startDateStr := c.QueryParam("start_date")
	endDateStr := c.QueryParam("end_date")

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid start_date format"})
	}
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid end_date format"})
	}

	ctx := context.Background()
	result, err := use.CalculateCustomerAndOrderStats(ctx, startDate, endDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to calculate customer/order stats"})
	}

	return c.JSON(http.StatusOK, entity.Response{
		Status:  200,
		Message: "Customer and Order Stats fetched successfully",
		Data:    result,
	})
}

func (a *AnalyticsLogs) GetProfitMarginByProduct(c echo.Context) error {
	use := usecase.AnalyticsLogs{MongoCon: a.MongoCon}

	startDateStr := c.QueryParam("start_date")
	endDateStr := c.QueryParam("end_date")

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid start_date format"})
	}
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid end_date format"})
	}

	ctx := context.Background()
	result, err := use.CalculateProfitMarginByProduct(ctx, startDate, endDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to calculate profit margin"})
	}

	return c.JSON(http.StatusOK, entity.Response{
		Status:  200,
		Message: "Profit margin by product fetched successfully",
		Data:    result,
	})
}
