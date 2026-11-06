package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthRoute(t *testing.T) {

	request := httptest.NewRequest("GET", "/health", nil)
	response := httptest.NewRecorder()

	handler := newHandler()
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf(
			"unexpected status: got %d, want %d",
			response.Code,
			http.StatusOK,
		)
	}

	gotContentType := response.Header().Get("Content-Type")
	wantContentType := "application/json"

	if gotContentType != wantContentType {
		t.Fatalf(
			"unexpected content type: got %q, want %q",
			gotContentType,
			wantContentType,
		)
	}

	gotBody := strings.TrimSpace(response.Body.String())
	wantBody := `{"status":"ok","service":"context-bridge"}`

	if gotBody != wantBody {
		t.Fatalf(
			"unexpected response body: got %q, want %q",
			gotBody,
			wantBody,
		)
	}
}
