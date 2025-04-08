package entity

import "time"

type Sale struct {
	OrderID         string    `bson:"order_id" json:"order_id"`
	ProductID       string    `bson:"product_id" json:"product_id"`
	CustomerID      string    `bson:"customer_id" json:"customer_id"`
	ProductName     string    `bson:"product_name" json:"product_name"`
	Category        string    `bson:"category" json:"category"`
	Region          string    `bson:"region" json:"region"`
	DateOfSale      time.Time `bson:"date_of_sale" json:"date_of_sale"`
	QuantitySold    int       `bson:"quantity_sold" json:"quantity_sold"`
	UnitPrice       float64   `bson:"unit_price" json:"unit_price"`
	Discount        float64   `bson:"discount" json:"discount"`
	ShippingCost    float64   `bson:"shipping_cost" json:"shipping_cost"`
	PaymentMethod   string    `bson:"payment_method" json:"payment_method"`
	CustomerName    string    `bson:"customer_name" json:"customer_name"`
	CustomerEmail   string    `bson:"customer_email" json:"customer_email"`
	CustomerAddress string    `bson:"customer_address" json:"customer_address"`
}

type RevenueGroupResult struct {
	Key   string  `bson:"_id" json:"key"`
	Total float64 `bson:"total" json:"total"`
}

type TotalRevenueResult struct {
	Total float64 `bson:"total" json:"total_revenue"`
}
