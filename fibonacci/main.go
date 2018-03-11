package main

import (
	"log"
	"net/http"
	"strconv"
)

func fibonacci(n uint64) uint64 {
	if n == 0 {
		return 0
	} else if n == 1 {
		return 1
	} else {
		return fibonacci(n-1) + fibonacci(n-2)
	}
}

func fibonacciHandler(w http.ResponseWriter, r *http.Request) {
	n, err := strconv.ParseUint(r.URL.Query().Get("broj"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Parametar broj mora biti prirodni broj."))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.FormatUint(fibonacci(n), 10)))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	http.HandleFunc("/", fibonacciHandler)
	http.HandleFunc("/health", healthHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
