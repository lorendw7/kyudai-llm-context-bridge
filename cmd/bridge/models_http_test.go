package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestListModelsHandlerAuthenticatedResponse(t *testing.T) {
	auth, err := NewTokenAuthenticator(map[string]Principal{
		"synthetic-token-a": {LabID: "lab-a", UserID: "user-1"},
		"synthetic-token-b": {LabID: "lab-b", UserID: "user-2"},
	})
	if err != nil {
		t.Fatal(err)
	}
	configured, err := NewModelAllowlist([]string{"local-model-b", "local-model-a"})
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range []struct {
		name      string
		allowlist ModelAllowlist
		want      string
	}{
		{"configured", configured, `{"models":["local-model-a","local-model-b"]}`},
		{"zero value", ModelAllowlist{}, `{"models":[]}`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.Handle("GET /v1/models", auth.RequireAuth(newListModelsHandler(tt.allowlist)))
			for _, token := range []string{"synthetic-token-a", "synthetic-token-b"} {
				r := httptest.NewRequest("GET", "/v1/models?model_id=unconfigured", nil)
				r.Header.Set("Authorization", "Bearer "+token)
				w := httptest.NewRecorder()
				mux.ServeHTTP(w, r)
				response := w.Result()
				response.Body.Close()
				if response.StatusCode != 200 || response.Header.Get("Content-Type") != "application/json" || response.Header.Get("Location") != "" || strings.TrimSpace(w.Body.String()) != tt.want {
					t.Errorf("unexpected models response: %d %v %s", response.StatusCode, response.Header, w.Body.String())
				}
			}
		})
	}
}

func TestListModelsHandlerRejectsInvalidAuthentication(t *testing.T) {
	auth, err := NewTokenAuthenticator(map[string]Principal{
		"synthetic-token-a": {LabID: "lab-a", UserID: "user-1"},
	})
	if err != nil {
		t.Fatal(err)
	}
	allowlist, err := NewModelAllowlist([]string{"synthetic-private-model"})
	if err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	mux.Handle("GET /v1/models", auth.RequireAuth(newListModelsHandler(allowlist)))
	for _, headers := range [][]string{nil, {"Bearer synthetic-unknown"}, {"Basic synthetic-token-a"}, {"Bearer synthetic-token-a", "Bearer synthetic-token-a"}} {
		r := httptest.NewRequest("GET", "/v1/models", nil)
		for _, value := range headers {
			r.Header.Add("Authorization", value)
		}
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		response := w.Result()
		response.Body.Close()
		if response.StatusCode != 401 || response.Header.Get("WWW-Authenticate") != "Bearer" || response.Header.Get("Content-Type") != "application/json" || strings.TrimSpace(w.Body.String()) != `{"code":"unauthorized"}` {
			t.Errorf("invalid authentication returned %d %v %s", response.StatusCode, response.Header, w.Body.String())
		}
	}
}
