package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Cotation struct {
	Usdbrl Usdbrl `json:"USDBRL"`
}

type Usdbrl struct {
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
	cotation, err := FindCotation()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cotation)
}

func FindCotation() (*Cotation, error) {
	req, err := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer a requisição: %v\n", err)
	}
	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao ler a resposta: %v\n", err)
	}
	var cotation Cotation
	err = json.Unmarshal(body, &cotation)
	if err != nil {
		return nil, err
	}
	return &cotation, nil
}
