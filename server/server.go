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

type Cambio struct {
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
		w.Write([]byte("Erro ao buscar as cotações."))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cambio)
}

func FindCotation() (*Cambio, error) {

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
	var cambio Cambio
	err = json.Unmarshal(body, &cambio)
	if err != nil {
		return nil, err
	}

	InsertCotation(ctxDatabase, cambio)

	return &cambio, nil
}

func InsertCotation(context context.Context, cambio Cambio) {
	db, err := CriaConexao()
	if err != nil {
		panic(err)
	}
	db.Create(&Cotation{
		Code:       cambio.Cotation.Code,
		Codein:     cambio.Cotation.Codein,
		Name:       cambio.Cotation.Name,
		High:       cambio.Cotation.High,
		Low:        cambio.Cotation.Low,
		VarBid:     cambio.Cotation.VarBid,
		PctChange:  cambio.Cotation.PctChange,
		Bid:        cambio.Cotation.Bid,
		Ask:        cambio.Cotation.Ask,
		Timestamp:  cambio.Cotation.Timestamp,
		CreateDate: cambio.Cotation.CreateDate,
	}).WithContext(context)
}

func CriaConexao() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(database), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&Cotation{})
	return db, nil
}
