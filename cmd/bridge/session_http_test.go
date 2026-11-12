package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func sessionCreationTestRouter(t *testing.T, store *MemorySessionStore) http.Handler {
	t.Helper()
	// Synthetic credentials only; this mux does not open a listener.
	auth, err := NewTokenAuthenticator(map[string]Principal{
		"synthetic-token-a": {LabID: "lab-a", UserID: "user-1"},
		"synthetic-token-b": {LabID: "lab-a", UserID: "user-2"},
		"synthetic-token-c": {LabID: "lab-b", UserID: "user-1"},
	})
	if err != nil {
		t.Fatal(err)
	}
	allowlist, err := NewModelAllowlist([]string{"synthetic-model-a"})
	if err != nil {
		t.Fatal(err)
	}
	return newAPIHandler(auth, store, allowlist)
}

func sessionCountForTest(store *MemorySessionStore) int {
	store.mu.Lock()
	defer store.mu.Unlock()
	return len(store.sessions)
}

func TestCreateSessionHandlerReturnsSavedOwnerBoundSession(t *testing.T) {
	store := NewMemorySessionStore()
	router := sessionCreationTestRouter(t, store)
	body := `{"schema_version":"1","policy_id":"full-history-v1"}`
	owners := []struct {
		token     string
		principal Principal
	}{
		{"synthetic-token-a", Principal{LabID: "lab-a", UserID: "user-1"}},
		{"synthetic-token-b", Principal{LabID: "lab-a", UserID: "user-2"}},
		{"synthetic-token-c", Principal{LabID: "lab-b", UserID: "user-1"}},
	}
	seen := make(map[string]bool)
	for _, owner := range owners {
		t.Run(owner.principal.LabID+"/"+owner.principal.UserID, func(t *testing.T) {
			r := httptest.NewRequest("POST", "/v1/sessions?user_id=forged&lab_id=forged", strings.NewReader(body))
			r.Header.Set("Authorization", "Bearer "+owner.token)
			r.Header.Set("X-User-ID", "forged")
			r.Header.Set("X-Lab-ID", "forged")
			w := httptest.NewRecorder()
			before := time.Now().UTC()
			router.ServeHTTP(w, r)
			after := time.Now().UTC()
			response := w.Result()
			defer response.Body.Close()
			if response.StatusCode != http.StatusCreated {
				t.Fatalf("status = %d, want 201", response.StatusCode)
			}
			// Result observes committed response headers, not later mutations.
			if got := response.Header.Get("Content-Type"); got != "application/json" {
				t.Errorf("committed Content-Type = %q, want application/json", got)
			}
			var result CreateSessionResponse
			if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
				t.Fatalf("invalid response JSON: %v", err)
			}
			if got := response.Header.Get("Location"); got != "/v1/sessions/"+result.ID {
				t.Errorf("committed Location = %q, want /v1/sessions/%s", got, result.ID)
			}
			decodedID, err := hex.DecodeString(result.ID)
			if err != nil || len(decodedID) != 16 || seen[result.ID] {
				t.Fatalf("invalid or repeated generated session ID: %q", result.ID)
			}
			seen[result.ID] = true
			saved, exists := store.Get(result.ID)
			if !exists || !saved.OwnedBy(owner.principal.LabID, owner.principal.UserID) {
				t.Fatalf("saved session missing or wrong owner: %+v", saved)
			}
			if result.PolicyID != "full-history-v1" || result.Revision != 1 || saved.PolicyID != result.PolicyID || saved.Revision != result.Revision || !saved.CreatedAt.Equal(result.CreatedAt) {
				t.Fatalf("response differs from saved session: %+v / %+v", result, saved)
			}
			if result.CreatedAt.Before(before) || result.CreatedAt.After(after) {
				t.Error("creation time was not generated during the request")
			}
			var fields map[string]json.RawMessage
			if err := json.Unmarshal(w.Body.Bytes(), &fields); err != nil || len(fields) != 4 {
				t.Error("response contains unexpected fields")
			}
		})
	}
	if got := sessionCountForTest(store); got != len(owners) {
		t.Errorf("saved sessions = %d, want %d", got, len(owners))
	}
}

func TestCreateSessionHandlerRejectsWithoutChangingStore(t *testing.T) {
	valid := `{"schema_version":"1","policy_id":"full-history-v1"}`
	tests := []struct {
		name   string
		body   string
		auth   []string
		status int
		code   string
	}{
		{"missing authentication", valid, nil, 401, "unauthorized"},
		{"unknown token", valid, []string{"Bearer synthetic-unknown"}, 401, "unauthorized"},
		{"duplicate authorization", valid, []string{"Bearer synthetic-token-a", "Bearer synthetic-token-a"}, 401, "unauthorized"},
		{"authentication before parsing", "broken", nil, 401, "unauthorized"},
	}
	invalidBodies := map[string]string{
		"empty body":          "",
		"malformed JSON":      `{"schema_version":`,
		"null":                `null`,
		"missing fields":      `{}`,
		"unsupported version": `{"schema_version":"99","policy_id":"full-history-v1"}`,
		"unsupported policy":  `{"schema_version":"1","policy_id":"unknown"}`,
		"wrong field type":    `{"schema_version":1,"policy_id":"full-history-v1"}`,
		"forged user":         `{"schema_version":"1","policy_id":"full-history-v1","user_id":"forged"}`,
		"forged lab":          `{"schema_version":"1","policy_id":"full-history-v1","lab_id":"forged"}`,
		"client session ID":   `{"schema_version":"1","policy_id":"full-history-v1","id":"forged"}`,
		"second object":       valid + `{}`,
		"trailing junk":       valid + `broken`,
	}
	for name, body := range invalidBodies {
		tests = append(tests, struct {
			name   string
			body   string
			auth   []string
			status int
			code   string
		}{name, body, []string{"Bearer synthetic-token-a"}, 422, "invalid_request"})
	}
	for name, body := range map[string]string{
		"oversized first value":         `{"schema_version":"1","policy_id":"` + strings.Repeat("x", 4096) + `"}`,
		"oversized trailing whitespace": valid + strings.Repeat(" ", 4097-len(valid)),
	} {
		tests = append(tests, struct {
			name   string
			body   string
			auth   []string
			status int
			code   string
		}{name, body, []string{"Bearer synthetic-token-a"}, 413, "request_too_large"})
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewMemorySessionStore()
			original := syntheticSession(t, "synthetic-existing-session")
			if err := store.Add(original); err != nil {
				t.Fatal(err)
			}
			router := sessionCreationTestRouter(t, store)
			r := httptest.NewRequest("POST", "/v1/sessions", strings.NewReader(tt.body))
			for _, value := range tt.auth {
				r.Header.Add("Authorization", value)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, r)
			response := w.Result()
			defer response.Body.Close()
			if response.StatusCode != tt.status || response.Header.Get("Content-Type") != "application/json" {
				t.Errorf("unexpected error response: %d %v", response.StatusCode, response.Header)
			}
			if response.Header.Get("Location") != "" {
				t.Error("failed creation returned Location")
			}
			var result map[string]string
			if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil || len(result) != 1 || result["code"] != tt.code {
				t.Errorf("error body = %q, want code %q", w.Body.String(), tt.code)
			}
			if sessionCountForTest(store) != 1 {
				t.Error("failed request added a session")
			}
			if got, exists := store.Get(original.ID); !exists || got != original {
				t.Error("failed request changed an existing session")
			}
		})
	}
}

func TestCreateSessionHandlerAcceptsExactBodyLimit(t *testing.T) {
	valid := `{"schema_version":"1","policy_id":"full-history-v1"}`
	body := valid + strings.Repeat(" ", 4096-len(valid))
	store := NewMemorySessionStore()
	r := httptest.NewRequest("POST", "/v1/sessions", strings.NewReader(body))
	r.Header.Set("Authorization", "Bearer synthetic-token-a")
	w := httptest.NewRecorder()
	sessionCreationTestRouter(t, store).ServeHTTP(w, r)
	if w.Code != http.StatusCreated || sessionCountForTest(store) != 1 {
		t.Fatalf("exact 4096-byte request failed: %d %s", w.Code, w.Body.String())
	}
}

func TestDecodeCreateSessionRequestReadsSingleObject(t *testing.T) {
	body := `{"schema_version":"1","policy_id":"full-history-v1"}`
	for _, suffix := range []string{"", " \n\t"} {
		r := httptest.NewRequest("POST", "/v1/sessions", strings.NewReader(body+suffix))
		got, err := decodeCreateSessionRequest(r)
		want := CreateSessionRequest{SchemaVersion: "1", PolicyID: "full-history-v1"}
		if err != nil || got != want {
			t.Fatalf("decode = %+v, %v; want %+v, nil", got, err, want)
		}
	}
}

func TestDecodeCreateSessionRequestRejectsMalformedOrAdditionalInput(t *testing.T) {
	valid := `{"schema_version":"1","policy_id":"full-history-v1"}`
	tests := map[string]string{
		"empty":         "",
		"truncated":     `{"schema_version":`,
		"wrong type":    `{"schema_version":1,"policy_id":"full-history-v1"}`,
		"forged owner":  `{"schema_version":"1","policy_id":"full-history-v1","user_id":"synthetic-other-user"}`,
		"second object": valid + `{}`,
		"second null":   valid + `null`,
		"trailing junk": valid + `broken`,
	}
	for name, body := range tests {
		t.Run(name, func(t *testing.T) {
			r := httptest.NewRequest("POST", "/v1/sessions", strings.NewReader(body))
			got, err := decodeCreateSessionRequest(r)
			if err == nil || got != (CreateSessionRequest{}) {
				t.Fatalf("invalid input returned %+v, %v", got, err)
			}
		})
	}
}

func TestDecodeCreateSessionRequestPreservesTrailingReadError(t *testing.T) {
	valid := `{"schema_version":"1","policy_id":"full-history-v1"}`
	r := httptest.NewRequest("POST", "/v1/sessions", strings.NewReader(valid+strings.Repeat(" ", 64)))
	r.Body = http.MaxBytesReader(httptest.NewRecorder(), r.Body, int64(len(valid)+1))
	_, err := decodeCreateSessionRequest(r)
	var tooLarge *http.MaxBytesError
	if !errors.As(err, &tooLarge) {
		t.Fatalf("trailing body-limit error was lost: %v", err)
	}
}

func TestWriteSessionErrorWritesStatusAndJSON(t *testing.T) {
	w := httptest.NewRecorder()
	writeSessionError(w, http.StatusUnprocessableEntity, "invalid_request")
	response := w.Result()
	defer response.Body.Close()
	if response.StatusCode != 422 || response.Header.Get("Content-Type") != "application/json" {
		t.Fatalf("unexpected error response: %d, %v", response.StatusCode, response.Header)
	}
	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if len(body) != 1 || body["code"] != "invalid_request" {
		t.Fatalf("unexpected error body: %v", body)
	}
}
