package main

import (
	"context"
	"database/sql"
	"encoding/json"
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

	log.Printf("Server running on port 8080")
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
	quotation, err := api_quotation_request()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}

	err = createQuotation(r.Context(), quotation.Quotation)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ApiResponse{Bid: quotation.Quotation.Bid})
}

func api_quotation_request() (USDBRL, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)

	if err != nil {
		return USDBRL{}, err
	}

	response, err := http.DefaultClient.Do(request)
	defer response.Body.Close()

	if err != nil {
		return USDBRL{}, err
	}

	var api_response USDBRL

	err = json.NewDecoder(response.Body).Decode(&api_response)

	if err != nil {
		return USDBRL{}, err
	}

	return api_response, nil
}

func dbConnection() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "file:goexpert.db")

	if err != nil {
		return nil, err
	}

	return db, nil
}

func createDataBase() error {
	db, err := dbConnection()

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

func createQuotation(ctx context.Context, quotation Quotation) error {
	db, err := dbConnection()

	if err != nil {
		return err
	}

	defer db.Close()

	query, err := db.Prepare(`
		INSERT INTO quotation(
			code,
			code_in,
			name,
			high,
			low,
			var_bid,
			pct_change,
			bid,
			ask,
			timestamp,
			create_date
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)

	if err != nil {
		return err
	}

	databaseContext, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()

	_, err = query.ExecContext(
		databaseContext,
		quotation.Code,
		quotation.CodeIn,
		quotation.Name,
		quotation.High,
		quotation.Low,
		quotation.VarBid,
		quotation.PctChange,
		quotation.Bid,
		quotation.Ask,
		quotation.Timestamp,
		quotation.CreateDate,
	)

	if err != nil {
		return err
	}

	return nil
}

type USDBRL struct {
	Quotation Quotation `json:"USDBRL"`
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

type ApiResponse struct {
	Bid string `json:"bid"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
