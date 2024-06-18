package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	bid, err := api_request()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(bid)
}

func api_request() (Bid, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/quotation", nil)

	if err != nil {
		return Bid{}, err
	}

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return Bid{}, err
	}

	var bid Bid
	err = json.NewDecoder(response.Body).Decode(&bid)

	if err != nil {
		return Bid{}, err
	}

	return bid, nil
}

type Bid struct {
	Bid string `json:"bid"`
}
