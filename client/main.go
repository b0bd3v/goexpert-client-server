package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	bid, err := api_request()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(bid)

	err = saveInFile(bid)

	if err != nil {
		log.Fatal(err)
	}
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

func saveInFile(bid Bid) error {
	file, err := os.Create("cotacao.txt")

	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.WriteString("Dólar: " + bid.Bid + "\n")

	if err != nil {
		return err
	}

	return nil
}

type Bid struct {
	Bid string `json:"bid"`
}
