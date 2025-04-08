package usecase

import (
	"analytics/entity"
	"analytics/loader"
	mongodb "analytics/repository"
	"analytics/utils"
	"context"
	"encoding/csv"
	"io"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AnalyticsLogs struct {
	MongoCon *mongodb.MongoCon
}

func (a *AnalyticsLogs) CalculateRevenueByType(ctx context.Context, startDate, endDate time.Time, aggType string) (interface{}, error) {
	collection := a.MongoCon.GetCollection("analytics", "analyticsLogs")

	matchStage := bson.M{
		"$match": bson.M{
			"date": bson.M{
				"$gte": startDate,
				"$lte": endDate,
			},
		},
	}

	var groupStage bson.M

	switch aggType {
	case "product":
		groupStage = bson.M{
			"$group": bson.M{
				"_id":   "$product_name",
				"total": bson.M{"$sum": "$revenue"},
			},
		}
	case "category":
		groupStage = bson.M{
			"$group": bson.M{
				"_id":   "$category",
				"total": bson.M{"$sum": "$revenue"},
			},
		}
	case "region":
		groupStage = bson.M{
			"$group": bson.M{
				"_id":   "$region",
				"total": bson.M{"$sum": "$revenue"},
			},
		}
	default: // total
		groupStage = bson.M{
			"$group": bson.M{
				"_id":   nil,
				"total": bson.M{"$sum": "$revenue"},
			},
		}
	}

	pipeline := []bson.M{matchStage, groupStage}
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Handle total revenue case
	if aggType == "" {
		var result entity.TotalRevenueResult
		if cursor.Next(ctx) {
			if err := cursor.Decode(&result); err != nil {
				return nil, err
			}
		}
		return result, nil
	}

	// Handle group results
	result := make(map[string]float64)
	for cursor.Next(ctx) {
		var row entity.RevenueGroupResult
		if err := cursor.Decode(&row); err != nil {
			return nil, err
		}
		result[row.Key] = row.Total
	}

	return result, nil
}

func (a *AnalyticsLogs) CalculateCustomerAndOrderStats(ctx context.Context, startDate, endDate time.Time) (interface{}, error) {
	collection := a.MongoCon.GetCollection("analytics", "analyticsLogs")

	matchStage := bson.M{"$match": bson.M{"date": bson.M{"$gte": startDate, "$lte": endDate}}}

	groupStage := bson.M{
		"$group": bson.M{
			"_id":             nil,
			"totalRevenue":    bson.M{"$sum": "$revenue"},
			"totalOrders":     bson.M{"$sum": 1},
			"uniqueCustomers": bson.M{"$addToSet": "$customer_id"},
		},
	}

	projectStage := bson.M{
		"$project": bson.M{
			"totalRevenue":   1,
			"totalOrders":    1,
			"totalCustomers": bson.M{"$size": "$uniqueCustomers"},
			"averageOrderValue": bson.M{"$cond": []interface{}{
				bson.M{"$eq": []interface{}{"$totalOrders", 0}}, 0,
				bson.M{"$divide": []interface{}{"$totalRevenue", "$totalOrders"}},
			}},
		},
	}

	cursor, err := collection.Aggregate(ctx, []bson.M{matchStage, groupStage, projectStage})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result map[string]interface{}
	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (a *AnalyticsLogs) CalculateProfitMarginByProduct(ctx context.Context, startDate, endDate time.Time) (interface{}, error) {
	collection := a.MongoCon.GetCollection("analytics", "analyticsLogs")

	matchStage := bson.M{"$match": bson.M{"date": bson.M{"$gte": startDate, "$lte": endDate}}}

	groupStage := bson.M{
		"$group": bson.M{
			"_id":          "$product_name",
			"totalRevenue": bson.M{"$sum": "$revenue"},
			"totalCost":    bson.M{"$sum": "$cost"},
		},
	}

	projectStage := bson.M{
		"$project": bson.M{
			"totalRevenue": 1,
			"totalCost":    1,
			"profitMargin": bson.M{"$cond": []interface{}{
				bson.M{"$eq": []interface{}{"$totalRevenue", 0}}, 0,
				bson.M{"$divide": []interface{}{bson.M{"$subtract": []interface{}{"$totalRevenue", "$totalCost"}}, "$totalRevenue"}},
			}},
		},
	}

	cursor, err := collection.Aggregate(ctx, []bson.M{matchStage, groupStage, projectStage})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	result := make(map[string]interface{})
	for cursor.Next(ctx) {
		var row struct {
			ID           string  `bson:"_id"`
			TotalRevenue float64 `bson:"totalRevenue"`
			TotalCost    float64 `bson:"totalCost"`
			ProfitMargin float64 `bson:"profitMargin"`
		}
		if err := cursor.Decode(&row); err != nil {
			return nil, err
		}
		result[row.ID] = row
	}

	return result, nil
}

func (a *AnalyticsLogs) RefreshAnalyticsData(ctx context.Context, csvPath string) error {
	collection := a.MongoCon.GetCollection("analytics", "analyticsLogs")

	// Optional: clear existing data
	if _, err := collection.DeleteMany(ctx, bson.M{}); err != nil {
		return err
	}

	return loader.LoadCSVToMongo(ctx, collection, csvPath)
}

func (a *AnalyticsLogs) ProcessUploadedCSV(ctx context.Context, reader io.Reader) error {
	csvReader := csv.NewReader(reader)

	// Read header row
	headers, err := csvReader.Read()
	if err != nil {
		return err
	}

	var sales []entity.Sale

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		sale, err := utils.MapCSVRecordToSale(headers, record)
		if err != nil {
			return err
		}

		sales = append(sales, sale)
	}

	// Bulk insert to MongoDB
	return a.BulkInsertSales(ctx, sales)
}

func (a *AnalyticsLogs) BulkInsertSales(ctx context.Context, sales []entity.Sale) error {
	collection := a.MongoCon.GetCollection("analytics", "analyticsLogs")

	if len(sales) == 0 {
		return nil
	}

	var docs []interface{}
	for _, sale := range sales {
		docs = append(docs, sale)
	}

	// Optional: Allow unordered insert
	opts := options.InsertMany().SetOrdered(false)

	_, err := collection.InsertMany(ctx, docs, opts)
	return err
}
