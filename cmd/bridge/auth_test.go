package main

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewPrincipalCreatesIdentity(t *testing.T) {
	principal, err := NewPrincipal("lab-a", "user-1")
	if err != nil {
		t.Fatalf("NewPrincipal returned an unexpected error: %v", err)
	}

	if principal.LabID != "lab-a" {
		t.Errorf("unexpected lab ID: got %q, want %q", principal.LabID, "lab-a")
	}
	if principal.UserID != "user-1" {
		t.Errorf("unexpected user ID: got %q, want %q", principal.UserID, "user-1")
	}
}

// All tokens and identities below are synthetic test fixtures.
func TestTokenAuthenticatorResolvesConfiguredIdentity(t *testing.T) {
	want := Principal{LabID: "lab-a", UserID: "user-1"}
	auth, err := NewTokenAuthenticator(map[string]Principal{
		"synthetic-token-a": want,
		"synthetic-token-b": {LabID: "lab-b", UserID: "user-2"},
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, scheme := range []string{"Bearer", "bearer", "BEARER"} {
		t.Run(scheme, func(t *testing.T) {
			r := httptest.NewRequest("POST", "/v1/sessions?user_id=user-2", strings.NewReader(`{"lab_id":"lab-b","user_id":"user-2"}`))
			r.Header.Set("Authorization", scheme+" synthetic-token-a")
			r.Header.Set("X-User-ID", "user-2")
			r.Header.Set("X-Lab-ID", "lab-b")
			got, err := auth.Authenticate(r)
			if err != nil || got != want {
				t.Fatalf("Authenticate = %+v, %v; want %+v, nil", got, err, want)
			}
		})
	}
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Bearer synthetic-token-b")
	got, err := auth.Authenticate(r)
	if err != nil || got != (Principal{LabID: "lab-b", UserID: "user-2"}) {
		t.Fatalf("second token resolved incorrectly: %+v, %v", got, err)
	}
}

func TestTokenAuthenticatorRejectsInvalidAuthorization(t *testing.T) {
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
		{"scheme only", []string{"Bearer"}},
		{"extra field", []string{"Bearer synthetic-token-a extra"}},
		{"wrong scheme", []string{"Basic synthetic-token-a"}},
		{"unknown token", []string{"Bearer synthetic-unknown-token"}},
		{"token case changed", []string{"Bearer SYNTHETIC-TOKEN-A"}},
		{"duplicate identical", []string{"Bearer synthetic-token-a", "Bearer synthetic-token-a"}},
		{"duplicate conflicting", []string{"Bearer synthetic-token-a", "Bearer synthetic-unknown-token"}},
		{"combined values", []string{"Bearer synthetic-token-a, Bearer synthetic-token-a"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", "/", nil)
			for _, value := range tt.values {
				r.Header.Add("Authorization", value)
			}
			got, err := auth.Authenticate(r)
			if err == nil || got != (Principal{}) {
				t.Fatalf("invalid request returned %+v, %v", got, err)
			}
			if strings.Contains(err.Error(), "synthetic-") {
				t.Fatal("authentication error exposed a credential")
			}
		})
	}
}

func TestNewTokenAuthenticatorRejectsInvalidConfiguration(t *testing.T) {
	valid := Principal{LabID: "lab-a", UserID: "user-1"}
	tests := []struct {
		name   string
		config map[string]Principal
	}{
		{"nil map", nil},
		{"empty map", map[string]Principal{}},
		{"empty token", map[string]Principal{"": valid}},
		{"leading space", map[string]Principal{" synthetic-token": valid}},
		{"embedded space", map[string]Principal{"synthetic token": valid}},
		{"tab", map[string]Principal{"synthetic\ttoken": valid}},
		{"newline", map[string]Principal{"synthetic-token\n": valid}},
		{"unicode space", map[string]Principal{"synthetic\u3000token": valid}},
		{"missing lab", map[string]Principal{"synthetic-token": {UserID: "user-1"}}},
		{"blank lab", map[string]Principal{"synthetic-token": {LabID: " \t", UserID: "user-1"}}},
		{"missing user", map[string]Principal{"synthetic-token": {LabID: "lab-a"}}},
		{"blank user", map[string]Principal{"synthetic-token": {LabID: "lab-a", UserID: " \t"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := NewTokenAuthenticator(tt.config); err == nil {
				t.Fatal("invalid configuration was accepted")
			}
		})
	}
}

func TestTokenAuthenticatorOwnsConfiguration(t *testing.T) {
	want := Principal{LabID: "lab-a", UserID: "user-1"}
	config := map[string]Principal{"synthetic-token-a": want}
	auth, err := NewTokenAuthenticator(config)
	if err != nil {
		t.Fatal(err)
	}
	config["synthetic-token-a"] = Principal{LabID: "lab-b", UserID: "user-2"}
	config["synthetic-token-b"] = want
	delete(config, "synthetic-token-a")
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Bearer synthetic-token-a")
	got, err := auth.Authenticate(r)
	if err != nil || got != want {
		t.Fatalf("caller mutation changed identity: %+v, %v", got, err)
	}
	r.Header.Set("Authorization", "Bearer synthetic-token-b")
	if _, err := auth.Authenticate(r); err == nil {
		t.Fatal("caller mutation added a credential")
	}
}

func TestZeroTokenAuthenticatorDeniesAccess(t *testing.T) {
	var auth TokenAuthenticator
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Bearer synthetic-token-a")
	if got, err := auth.Authenticate(r); err == nil || got != (Principal{}) {
		t.Fatalf("zero authenticator returned %+v, %v", got, err)
	}
}

func TestNewPrincipalRejectsMissingIdentity(t *testing.T) {
	tests := []struct {
		name   string
		labID  string
		userID string
	}{
		{name: "missing lab ID", userID: "user-1"},
		{name: "blank lab ID", labID: "   ", userID: "user-1"},
		{name: "missing user ID", labID: "lab-a"},
		{name: "blank user ID", labID: "lab-a", userID: "   "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := NewPrincipal(tt.labID, tt.userID); err == nil {
				t.Fatal("NewPrincipal returned nil error for invalid identity")
			}
		})
	}
}

func TestPrincipalCanAccessOnlyOwnedSession(t *testing.T) {
	session, err := NewSession(
		"session-1",
		"lab-a",
		"user-1",
		"full-history-v1",
		time.Date(2026, time.November, 10, 9, 0, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("NewSession returned an unexpected error: %v", err)
	}

	tests := []struct {
		name      string
		labID     string
		userID    string
		wantAllow bool
	}{
		{name: "owner", labID: "lab-a", userID: "user-1", wantAllow: true},
		{name: "same user in another lab", labID: "lab-b", userID: "user-1", wantAllow: false},
		{name: "another user in same lab", labID: "lab-a", userID: "user-2", wantAllow: false},
		{name: "another user and lab", labID: "lab-b", userID: "user-2", wantAllow: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			principal, err := NewPrincipal(tt.labID, tt.userID)
			if err != nil {
				t.Fatalf("NewPrincipal returned an unexpected error: %v", err)
			}

			if got := principal.CanAccessSession(session); got != tt.wantAllow {
				t.Errorf("CanAccessSession() = %v, want %v", got, tt.wantAllow)
			}
		})
	}
}
