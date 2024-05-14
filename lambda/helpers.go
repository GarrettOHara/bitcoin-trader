package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
)

func generateUUID() []byte {
	// Generate a new UUID
	uuidObj := uuid.New()

	// Convert the UUID to a byte slice
	idBytes := uuidObj[:]

	// Return the 16-byte UUID
	return idBytes
}

func printOrderConfig(order Order) {
	// Convert order struct to map[string]interface{}
	payloadMap := map[string]interface{}{
		"side": order.Side,
		"order_configuration": map[string]interface{}{
			"market_market_ioc": map[string]interface{}{
				"quote_size": order.OrderConfiguration.MarketMarketIOC.QuoteSize,
			},
		},
		"product_id":      order.ProductID,
		"client_order_id": fmt.Sprintf("%x", order.ClientOrderID),
	}

	// Marshal the map to JSON with indentation
	prettyPayload, err := json.MarshalIndent(payloadMap, "", "  ")
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		os.Exit(1)
	}

	// Print the JSON payload
	fmt.Println(string(prettyPayload))
}
