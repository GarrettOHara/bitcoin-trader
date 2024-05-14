package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/davecgh/go-spew/spew"
)

func LambdaHandler(request events.CloudWatchEvent) error {
	fmt.Printf("received event of type %q\n", request.DetailType)

	// Spew handles complex payloads/string objects with ease
	// https://pkg.go.dev/github.com/davecgh/go-spew/spew
	spew.Dump(request)

	const tradeAmount = 1
	url := "https://api.coinbase.com/api/v3/brokerage/orders"
	clientOrderId := generateUUID()
	requestJwt, err := generateJtw()
	if err != nil {
		fmt.Println("Error generating jwt: ", err)
		os.Exit(1)
	}

	orderConfig := OrderConfiguration{
		MarketMarketIOC: MarketMarketIOC{
			QuoteSize: fmt.Sprintf("%d", tradeAmount),
		},
	}

	order := Order{
		Side:               "BUY",
		OrderConfiguration: orderConfig,
		ProductID:          "BTC-USD",
		ClientOrderID:      clientOrderId,
	}

	printOrderConfig(order)

	// Convert the map to a JSON string
	payload, err := json.Marshal(order)
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
	// Read the response body and print it
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		os.Exit(1)
	}
	fmt.Println("Response Body:", string(body))

	spew.Dump(resp.Body)
	return nil
}

func main() {
	lambda.Start(LambdaHandler)
}
