package main

import (
	"log"
	"net/http"
	"strconv"
)

func fibonacci(n uint64) uint64 {
	if n == 0 {
		return 0
	}
	a := uint64(0)
	b := uint64(1)

	for n > 1 {
		tmp := a + b
		a = b
		b = tmp
		n--
	}
	return b
}

func fibonacciHandler(w http.ResponseWriter, r *http.Request) {
	n, err := strconv.ParseUint(r.URL.Query().Get("n"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Parametar n mora biti prirodni broj."))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.FormatUint(fibonacci(n), 10)))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
}

func main() {
	http.HandleFunc("/", homePageHandler)
	http.HandleFunc("/api/fibonacci", fibonacciHandler)
	http.HandleFunc("/health", healthHandler)

	log.Fatal(http.ListenAndServe(":8888", nil))
}
