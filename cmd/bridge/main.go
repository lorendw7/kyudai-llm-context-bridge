package main

import (
	"fmt"
	"log"
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

func newAPIHandler(
	auth TokenAuthenticator,
	store *MemorySessionStore,
	allowlist ModelAllowlist,
) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler)
	mux.Handle("POST /v1/sessions",
		auth.RequireAuth(newCreateSessionHandler(store)))

	mux.Handle("GET /v1/sessions",
		auth.RequireAuth(newListSessionsHandler(store)))

	mux.Handle("GET /v1/sessions/{id}",
		auth.RequireAuth(newGetSessionHandler(store)))

	mux.Handle("GET /v1/models",
		auth.RequireAuth(newListModelsHandler(allowlist)))

	return mux
}

func main() {
	auth, err := authFromEnvironment()
	if err != nil {
		log.Fatal("Invalid authentication configuration: ", err)
	}

	allowlist, err := modelsFromEnvironment()
	if err != nil {
		log.Fatal("Invalid model configuration: ", err)
	}

	store := NewMemorySessionStore()
	handler := newAPIHandler(auth, store, allowlist)

	fmt.Println("Starting Bridge at http://127.0.0.1:8080")

	if err := http.ListenAndServe("127.0.0.1:8080", handler); err != nil {
		log.Fatal("Server stopped: ", err)
	}
}
