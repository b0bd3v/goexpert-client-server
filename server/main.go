package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

const cotacaoURL string = "https://economia.awesomeapi.com.br/json/last/USD-BRL"

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
	//log.Println("GET /cotacao")
	//
	//err := saveQuotationToDB()
	//
	//if err != nil {
	//	return err
	//}

	w.WriteHeader(http.StatusOK)
}

//func saveQuotationToDB() error {}

func createDataBase() error {
	db, err := sql.Open("sqlite3", "file:goexpert.db")

	if err != nil {
		return err
	}

	defer db.Close()

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS cotacao(
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
