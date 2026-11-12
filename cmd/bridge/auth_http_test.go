package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequireAuthRejectsBeforeCallingHandler(t *testing.T) {
	// Synthetic credentials only; no listener or external service is used.
	auth, err := NewTokenAuthenticator(map[string]Principal{
		"synthetic-token-a": {LabID: "lab-a", UserID: "user-1"},
	})
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name   string
		values []string
	}{
		{"missing", nil},
		{"empty", []string{""}},
		{"missing token", []string{"Bearer"}},
		{"wrong scheme", []string{"Basic synthetic-token-a"}},
		{"unknown token", []string{"Bearer synthetic-unknown-token"}},
		{"token case changed", []string{"Bearer SYNTHETIC-TOKEN-A"}},
		{"extra field", []string{"Bearer synthetic-token-a extra"}},
		{"duplicate", []string{"Bearer synthetic-token-a", "Bearer synthetic-token-a"}},
		{"conflicting", []string{"Bearer synthetic-token-a", "Bearer synthetic-unknown-token"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			protected := auth.RequireAuth(func(w http.ResponseWriter, r *http.Request, p Principal) {
				calls++
				w.WriteHeader(http.StatusNoContent)
			})
			mux := http.NewServeMux()
			mux.Handle("POST /test/protected", protected)
			r := httptest.NewRequest("POST", "/test/protected", nil)
			for _, value := range tt.values {
				r.Header.Add("Authorization", value)
			}
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, r)
			if calls != 0 {
				t.Fatalf("unauthenticated request reached business handler %d times", calls)
			}
			if w.Code != http.StatusUnauthorized {
				t.Fatalf("status = %d, want 401", w.Code)
			}
			// Result snapshots the headers committed when WriteHeader was called.
			response := w.Result()
			defer response.Body.Close()
			if got := response.Header.Get("Content-Type"); got != "application/json" {
				t.Errorf("Content-Type = %q", got)
			}
			if got := response.Header.Get("WWW-Authenticate"); got != "Bearer" {
				t.Errorf("WWW-Authenticate = %q", got)
			}
			var body map[string]string
			if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
				t.Fatalf("invalid JSON error response: %v", err)
			}
			if len(body) != 1 || body["code"] != "unauthorized" {
				t.Errorf("unexpected learning-stage error response: %v", body)
			}
			if strings.Contains(w.Body.String(), "synthetic-") {
				t.Error("response leaked a credential")
			}
		})
	}
}

func TestRequireAuthPassesIdentityAndPreservesHandlerResponse(t *testing.T) {
	identities := map[string]Principal{
		"synthetic-token-a": {LabID: "lab-a", UserID: "user-1"},
		"synthetic-token-b": {LabID: "lab-b", UserID: "user-2"},
	}
	auth, err := NewTokenAuthenticator(identities)
	if err != nil {
		t.Fatal(err)
	}
	calls := 0
	var gotPrincipal Principal
	var gotRequest *http.Request
	protected := auth.RequireAuth(func(w http.ResponseWriter, r *http.Request, p Principal) {
		calls++
		gotPrincipal, gotRequest = p, r
		w.Header().Set("X-Test-Handler", "reached")
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte("synthetic-handler-response"))
	})
	if calls != 0 {
		t.Fatal("business handler executed during wrapper construction")
	}
	mux := http.NewServeMux()
	mux.Handle("POST /test/protected", protected)
	for token, want := range identities {
		r := httptest.NewRequest("POST", "/test/protected?user_id=forged", strings.NewReader(`{"user_id":"forged","lab_id":"forged"}`))
		r.Header.Set("Authorization", "bearer "+token)
		r.Header.Set("X-User-ID", "forged")
		w := httptest.NewRecorder()
		before := calls
		mux.ServeHTTP(w, r)
		if calls != before+1 || gotPrincipal != want || gotRequest != r {
			t.Fatalf("incorrect delegation: calls=%d, identity=%+v, same request=%v", calls-before, gotPrincipal, gotRequest == r)
		}
		if w.Code != http.StatusAccepted || w.Body.String() != "synthetic-handler-response" {
			t.Fatalf("middleware changed handler response: %d %q", w.Code, w.Body.String())
		}
		response := w.Result()
		response.Body.Close()
		if response.Header.Get("X-Test-Handler") != "reached" || response.Header.Get("WWW-Authenticate") != "" {
			t.Fatal("unexpected success response headers")
		}
	}
}
