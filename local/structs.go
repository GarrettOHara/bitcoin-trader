package main

type OrderConfiguration struct {
	MarketMarketIOC MarketMarketIOC `json:"market_market_ioc"`
}

type MarketMarketIOC struct {
	QuoteSize string `json:"quote_size"`
}

type Order struct {
	Side               string             `json:"side"`
	OrderConfiguration OrderConfiguration `json:"order_configuration"`
	ProductID          string             `json:"product_id"`
	ClientOrderID      []byte             `json:"client_order_id"`
}
