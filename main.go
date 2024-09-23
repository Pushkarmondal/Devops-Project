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
    if err = tmpl.Execute(w, nil); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func handleConvert(w http.ResponseWriter, r *http.Request) {
    amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
    if err != nil {
        http.Error(w, "Invalid amount", http.StatusBadRequest)
        return
    }
    from := r.FormValue("from")
    to := r.FormValue("to")

    rate, err := fetchExchangeRate(from, to)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    result := ConversionResult{
        Amount: amount,
        From:   from,
        To:     to,
        Result: amount * rate,
    }

    w.Header().Set("Content-Type", "application/json")
    if err = json.NewEncoder(w).Encode(result); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func fetchExchangeRate(from, to string) (float64, error) {
    apiKey := "b0ec366a979f3e3fecf7348c13251214" 
    url := fmt.Sprintf("http://api.exchangeratesapi.io/v1/latest?access_key=%s&base=%s&symbols=%s", apiKey, from, to)

    resp, err := http.Get(url)
    if err != nil {
        return 0, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return 0, err
    }

    var response ExchangeRateResponse
    if err = json.Unmarshal(body, &response); err != nil {
        return 0, err
    }

    if !response.Success {
        return 0, fmt.Errorf("failed to fetch exchange rates")
    }

    rate, exists := response.Rates[to]
    if !exists {
        return 0, fmt.Errorf("rate not found for currency: %s", to)
    }

    return rate, nil
}
