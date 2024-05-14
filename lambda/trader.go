package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/davecgh/go-spew/spew"
)

func LambdaHandler(request events.CloudWatchEvent) error {
	// Perhaps need %v or %s
	log.Printf("received event of type %q\n", request.DetailType)

	// Spew handles complex payloads/string objects with ease
	// https://pkg.go.dev/github.com/davecgh/go-spew/spew
	// Compare output of this to log and use better one
	spew.Dump(request)

	const tradeAmount = 1
	url := "https://api.coinbase.com/api/v3/brokerage/orders"
	clientOrderId := generateUUID()
	requestJwt, err := generateJtw()
	if err != nil {
		log.Printf("Error generating jwt: %v", err)
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
		log.Printf("Error encoding Request payload: %v", err)
		os.Exit(1)
	}

	// Create a new HTTP request with POST method
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		os.Exit(1)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+requestJwt)
	req.Header.Set("Content-Type", "application/json")

	// Perform the HTTP POST request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	log.Printf("Response Status: %d", resp.Status)

	// Read the response body and print it
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		os.Exit(1)
	}
	log.Printf("Response Body: %s", string(body))

	// Compare the output of this and log and use the better one
	spew.Dump(resp.Body)
	return nil
}

func main() {
	lambda.Start(LambdaHandler)
}
