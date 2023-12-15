package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const database string = "client-server-api.db"

type Exchange struct {
	Cotation Cotation `json:"USDBRL"`
}

type Cotation struct {
	ID         int    `gorm:"primarykey" json:"-"`
	Code       string `json:"code"`
	Codein     string `json:"codein"`
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

func main() {
	http.HandleFunc("/cotacao", FindCotationHandler)
	http.ListenAndServe(":8080", nil)
}

func FindCotationHandler(w http.ResponseWriter, r *http.Request) {
	cambio, err := FindCotation()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error to find cotation."))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cambio)
}

func FindCotation() (*Exchange, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	ctxDatabase, cancelDatabase := context.WithTimeout(context.Background(), 10*time.Millisecond)
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		panic(err)
	}
	defer cancel()
	defer cancelDatabase()
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		panic(err)
	}
	var exchange Exchange
	err = json.Unmarshal(body, &exchange)
	if err != nil {
		return nil, err
	}

	InsertCotation(ctxDatabase, exchange)

	return &exchange, nil
}

func InsertCotation(context context.Context, exchange Exchange) {
	db, err := CreateConection()
	if err != nil {
		panic(err)
	}
	if err := db.WithContext(context).Create(&Cotation{
		Code:       exchange.Cotation.Code,
		Codein:     exchange.Cotation.Codein,
		Name:       exchange.Cotation.Name,
		High:       exchange.Cotation.High,
		Low:        exchange.Cotation.Low,
		VarBid:     exchange.Cotation.VarBid,
		PctChange:  exchange.Cotation.PctChange,
		Bid:        exchange.Cotation.Bid,
		Ask:        exchange.Cotation.Ask,
		Timestamp:  exchange.Cotation.Timestamp,
		CreateDate: exchange.Cotation.CreateDate,
	}).Error; err != nil {
		panic(err.Error())
	}

}

func CreateConection() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(database), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&Cotation{})
	return db, nil
}
