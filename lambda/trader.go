package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/davecgh/go-spew/spew"
)

func LambdaHandler(request events.CloudWatchEvent) error {
	fmt.Printf("received event of type %q\n", request.DetailType)

	// Spew handles complex payloads/string objects with ease
	// https://pkg.go.dev/github.com/davecgh/go-spew/spew
	spew.Dump(request)

	url := "https://api.coinbase.com/api/v3/brokerage/orders"
	method := "POST"

	const tradeAmount = 1
	const clientOrderId = "xxxx-xxxx-xxxx-xxxx"

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

	// Convert the map to a JSON string
	payload, err := json.Marshal(payloadMap)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return err
	}

	payloadReader := strings.NewReader(string(payload))

	// Instantiate http client and construct POST Request
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payloadReader)

	if err != nil {
		fmt.Println(err)
		return err
	}

	// Add HTTP request headers
	req.Header.Add("Content-Type", "application/json")

	// Send HTTP request
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	// Close HTTP client when function ends
	defer res.Body.Close()

	spew.Dump(res.Body)
	return nil
}

func main() {
	lambda.Start(LambdaHandler)
}
