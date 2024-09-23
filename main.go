// main.go
package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
)

type ExchangeRateResponse struct {
	Success   bool               `json:"success"`
	Timestamp int                `json:"timestamp"`
	Base      string             `json:"base"`
	Date      string             `json:"date"`
	Rates     map[string]float64 `json:"rates"`
}

type ConversionResult struct {
	Amount float64 `json:"amount"`
	From   string  `json:"from"`
	To     string  `json:"to"`
	Result float64 `json:"result"`
}

func main() {
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/convert", handleConvert)

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func handleConvert(w http.ResponseWriter, r *http.Request) {
	amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
	if err != nil {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}
	from := r.FormValue("from")
	to := r.FormValue("to")

	url := fmt.Sprintf("http://api.exchangeratesapi.io/v1/latest?access_key=b0ec366a979f3e3fecf7348c13251214&base=%s&symbols=%s", from, to)

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var rateResponse ExchangeRateResponse
	err = json.Unmarshal(body, &rateResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !rateResponse.Success {
		http.Error(w, "Failed to fetch exchange rates", http.StatusInternalServerError)
		return
	}

	rate := rateResponse.Rates[to]
	convertedAmount := amount * rate

	result := ConversionResult{
		Amount: amount,
		From:   from,
		To:     to,
		Result: convertedAmount,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
