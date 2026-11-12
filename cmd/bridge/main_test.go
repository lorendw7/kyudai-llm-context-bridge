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

func TestAPIHandlerRoutesAndAuthentication(t *testing.T) {
	auth, err := NewTokenAuthenticator(map[string]Principal{
		"synthetic-api-token": {LabID: "lab-a", UserID: "user-1"},
	})
	if err != nil {
		t.Fatal(err)
	}
	allowlist, err := NewModelAllowlist([]string{"synthetic-model-b", "synthetic-model-a"})
	if err != nil {
		t.Fatal(err)
	}
	store := NewMemorySessionStore()
	handler := newAPIHandler(auth, store, allowlist)
	tests := []struct {
		name, method, path, token string
		status                    int
		body                      string
	}{
		{"public health", "GET", "/health", "", 200, `{"status":"ok","service":"context-bridge"}`},
		{"models", "GET", "/v1/models", "synthetic-api-token", 200, `{"models":["synthetic-model-a","synthetic-model-b"]}`},
		{"empty sessions", "GET", "/v1/sessions", "synthetic-api-token", 200, `{"sessions":[]}`},
		{"protect models", "GET", "/v1/models", "", 401, `{"code":"unauthorized"}`},
		{"protect list", "GET", "/v1/sessions", "", 401, `{"code":"unauthorized"}`},
		{"protect creation", "POST", "/v1/sessions", "", 401, `{"code":"unauthorized"}`},
		{"protect retrieval", "GET", "/v1/sessions/synthetic-absent", "", 401, `{"code":"unauthorized"}`},
		{"unknown session", "GET", "/v1/sessions/synthetic-absent", "synthetic-api-token", 404, `{"code":"session_not_found"}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.method, tt.path, nil)
			if tt.token != "" {
				r.Header.Set("Authorization", "Bearer "+tt.token)
			}
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, r)
			response := w.Result()
			defer response.Body.Close()
			if response.StatusCode != tt.status || response.Header.Get("Content-Type") != "application/json" || strings.TrimSpace(w.Body.String()) != tt.body {
				t.Errorf("unexpected response: %d %v %s", response.StatusCode, response.Header, w.Body.String())
			}
		})
	}
	for _, path := range []string{"/v1/models", "/v1/sessions", "/v1/sessions/synthetic-absent"} {
		r := httptest.NewRequest("DELETE", path, nil)
		r.Header.Set("Authorization", "Bearer synthetic-api-token")
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("DELETE %s status = %d", path, w.Code)
		}
	}
	if sessionCountForTest(store) != 0 {
		t.Error("routing checks unexpectedly wrote sessions")
	}
}
