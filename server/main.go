package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"time"
)

func main() {
	var err error
	err = createDataBase()

	if err != nil {
		log.Fatal(err)
	}

	err = startServer()

	if err != nil {
		log.Fatal(err)
	}
}

func startServer() error {
	http.HandleFunc("/quotation", mainHandler)
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		return err
	}

	return nil
}

func mainHandler(w http.ResponseWriter, r *http.Request) {

	err := api_quotation_request()

	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusOK)
}

func api_quotation_request() error {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)

	if err != nil {
		return err
	}

	response, err := http.DefaultClient.Do(request)
	defer response.Body.Close()

	if err != nil {
		return err
	}

	var api_response map[string]Quotation

	err = json.NewDecoder(response.Body).Decode(&api_response)

	fmt.Printf("---> %+v\n", api_response)

	return nil
}

func createDataBase() error {
	db, err := sql.Open("sqlite3", "file:goexpert.db")

	if err != nil {
		return err
	}

	defer db.Close()

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS quotation(
		code TEXT, 
		code_in TEXT, 
		name TEXT, 
		high TEXT, 
		low TEXT,
		var_bid TEXT,
		pct_change TEXT,
		bid TEXT,
		ask TEXT,
		timestamp TEXT,
		create_date TEXT
	)`)

	if err != nil {
		return err
	}

	return nil
}

type Quotation struct {
	Code       string `json:"code"`
	CodeIn     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}
