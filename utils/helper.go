package utils

import (
	"analytics/entity"
	"strconv"
	"time"
)

// Helper to convert a CSV record to entity.Sale

func MapCSVRecordToSale(headers []string, record []string) (entity.Sale, error) {
	layout := "2006-01-02"

	dateOfSale, err := time.Parse(layout, record[6])
	if err != nil {
		return entity.Sale{}, err
	}

	quantity, err := strconv.Atoi(record[7])
	if err != nil {
		return entity.Sale{}, err
	}

	unitPrice, err := strconv.ParseFloat(record[8], 64)
	if err != nil {
		return entity.Sale{}, err
	}

	discount, err := strconv.ParseFloat(record[9], 64)
	if err != nil {
		return entity.Sale{}, err
	}

	shippingCost, err := strconv.ParseFloat(record[10], 64)
	if err != nil {
		return entity.Sale{}, err
	}

	return entity.Sale{
		OrderID:         record[0],
		ProductID:       record[1],
		CustomerID:      record[2],
		ProductName:     record[3],
		Category:        record[4],
		Region:          record[5],
		DateOfSale:      dateOfSale,
		QuantitySold:    quantity,
		UnitPrice:       unitPrice,
		Discount:        discount,
		ShippingCost:    shippingCost,
		PaymentMethod:   record[11],
		CustomerName:    record[12],
		CustomerEmail:   record[13],
		CustomerAddress: record[14],
	}, nil
}
