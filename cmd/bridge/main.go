package main

import (
	"fmt"
	"net/http"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, `{"status":"ok","service":"context-bridge"}`)
}

func newHandler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler)

	return mux
}

func main() {
	handler := newHandler()

	fmt.Println("Bridge running at http://127.0.0.1:8080")

	err := http.ListenAndServe("127.0.0.1:8080", handler)

	if err != nil {
		fmt.Println("Server stopped", err)
	}
}
