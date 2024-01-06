package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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

func debugPrint(payloadMap map[string]interface{}) {
	prettyPayload, err := json.MarshalIndent(payloadMap, "", "  ")
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		os.Exit(1)
	}
	fmt.Println(string(prettyPayload))
}

// Order Documentation:
// https://docs.cloud.coinbase.com/advanced-trade-api/reference/retailbrokerageapi_postorder
func main() {

	const tradeAmount = 1

	url := "https://api.coinbase.com/api/v3/brokerage/orders"
	clientOrderId := generateUUID()
	requestJwt, err := generateJtw()
	if err != nil {
		fmt.Println("Error generating jwt: ", err)
		os.Exit(1)
	}

	payloadMap := map[string]interface{}{
		"side": "BUY",
		"order_configuration": map[string]interface{}{
			"market_market_ioc": map[string]interface{}{
				"quote_size": fmt.Sprintf("%d", tradeAmount),
			},
		},
		"product_id":      "BTC-USD",
		"client_order_id": clientOrderId,
	}
	debugPrint(payloadMap)

	// Convert the map to a JSON string
	payload, err := json.Marshal(payloadMap)
	if err != nil {
		fmt.Println("Error encoding Request payload:", err)
		os.Exit(1)
	}

	// Create a new HTTP request with POST method
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		os.Exit(1)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+requestJwt)
	req.Header.Set("Content-Type", "application/json")

	// Perform the HTTP POST request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Process the response as needed
	fmt.Println("Response Status:", resp.Status)
}
