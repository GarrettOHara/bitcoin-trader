package main

type Order struct {
	Side               string                 `json:"side"`
	OrderConfiguration map[string]interface{} `json:"order_configuration"`
	ProductID          string                 `json:"product_id"`
	ClientOrderID      string                 `json:"client_order_id"`
}
