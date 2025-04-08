package loader

import (
	"bufio"
	"context"
	"encoding/csv"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func LoadCSVToMongo(ctx context.Context, collection *mongo.Collection, csvPath string) error {
	file, err := os.Open(csvPath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	headers := records[0]
	records = records[1:]

	var docs []interface{}

	for _, row := range records {
		rowMap := map[string]string{}
		for i, val := range row {
			rowMap[headers[i]] = val
		}

		date, _ := time.Parse("2006-01-02", rowMap["Date of Sale"])
		quantity, _ := strconv.Atoi(rowMap["Quantity Sold"])
		unitPrice, _ := strconv.ParseFloat(rowMap["Unit Price"], 64)
		discount, _ := strconv.ParseFloat(rowMap["Discount"], 64)
		shippingCost, _ := strconv.ParseFloat(rowMap["Shipping Cost"], 64)

		revenue := float64(quantity) * unitPrice * (1 - discount)
		cost := float64(quantity) * unitPrice * 0.5 // assuming 50% cost basis

		doc := bson.M{
			"order_id":         rowMap["Order ID"],
			"product_id":       rowMap["Product ID"],
			"customer_id":      rowMap["Customer ID"],
			"product_name":     rowMap["Product Name"],
			"category":         rowMap["Category"],
			"region":           rowMap["Region"],
			"date":             date,
			"quantity_sold":    quantity,
			"unit_price":       unitPrice,
			"discount":         discount,
			"shipping_cost":    shippingCost,
			"payment_method":   rowMap["Payment Method"],
			"customer_name":    rowMap["Customer Name"],
			"customer_email":   rowMap["Customer Email"],
			"customer_address": rowMap["Customer Address"],
			"revenue":          revenue,
			"cost":             cost,
		}

		docs = append(docs, doc)
	}

	if len(docs) > 0 {
		_, err = collection.InsertMany(ctx, docs)
	}

	return err
}
