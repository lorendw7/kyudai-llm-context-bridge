package main

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestListSessionsHandlerReturnsOnlyOwnedMetadata(t *testing.T) {
	store := NewMemorySessionStore()
	owners := []struct {
		token string
		lab   string
		user  string
	}{
		{"synthetic-token-a", "lab-a", "user-1"},
		{"synthetic-token-b", "lab-a", "user-2"},
		{"synthetic-token-c", "lab-b", "user-1"},
	}
	want := make(map[string]map[string]Session)
	for i, owner := range owners {
		want[owner.token] = make(map[string]Session)
		for j := 0; j < 2; j++ {
			session := syntheticSession(t, fmt.Sprintf("synthetic-list-%d-%d", i, j))
			session.LabID, session.OwnerUserID = owner.lab, owner.user
			if err := store.Add(session); err != nil {
				t.Fatal(err)
			}
			want[owner.token][session.ID] = session
		}
	}
	router := sessionCreationTestRouter(t, store)
	for _, owner := range owners {
		t.Run(owner.lab+"/"+owner.user, func(t *testing.T) {
			r := httptest.NewRequest("GET", "/v1/sessions?lab_id=forged&user_id=forged", nil)
			r.Header.Set("Authorization", "Bearer "+owner.token)
			r.Header.Set("X-Lab-ID", "forged")
			r.Header.Set("X-User-ID", "forged")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, r)
			response := w.Result()
			defer response.Body.Close()
			if response.StatusCode != 200 || response.Header.Get("Content-Type") != "application/json" || response.Header.Get("Location") != "" {
				t.Fatalf("unexpected list response: %d %v", response.StatusCode, response.Header)
			}
			var result ListSessionsResponse
			if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
				t.Fatal(err)
			}
			if len(result.Sessions) != len(want[owner.token]) {
				t.Fatalf("list count = %d, want %d", len(result.Sessions), len(want[owner.token]))
			}
			seen := make(map[string]bool)
			for _, item := range result.Sessions {
				expected, exists := want[owner.token][item.ID]
				if !exists || seen[item.ID] {
					t.Fatalf("foreign or repeated session: %s", item.ID)
				}
				seen[item.ID] = true
				if item.PolicyID != expected.PolicyID || item.Revision != expected.Revision || !item.CreatedAt.Equal(expected.CreatedAt) {
					t.Errorf("metadata differs from storage: %+v / %+v", item, expected)
				}
			}
			var envelope map[string]json.RawMessage
			if err := json.Unmarshal(w.Body.Bytes(), &envelope); err != nil || len(envelope) != 1 {
				t.Fatal("unexpected response envelope")
			}
			var items []map[string]json.RawMessage
			if err := json.Unmarshal(envelope["sessions"], &items); err != nil {
				t.Fatal(err)
			}
			for _, item := range items {
				if len(item) != 4 {
					t.Error("list exposed unexpected session fields")
				}
			}
		})
	}
	if sessionCountForTest(store) != 6 {
		t.Error("listing changed session count")
	}
	for _, entries := range want {
		for id, expected := range entries {
			if got, exists := store.Get(id); !exists || got != expected {
				t.Errorf("listing changed session %s", id)
			}
		}
	}
}

func TestListSessionsHandlerReturnsEmptyArray(t *testing.T) {
	for _, foreignPresent := range []bool{false, true} {
		t.Run(fmt.Sprintf("foreign_present_%v", foreignPresent), func(t *testing.T) {
			store := NewMemorySessionStore()
			if foreignPresent {
				if err := store.Add(syntheticSession(t, "synthetic-other-owner")); err != nil {
					t.Fatal(err)
				}
			}
			r := httptest.NewRequest("GET", "/v1/sessions", nil)
			r.Header.Set("Authorization", "Bearer synthetic-token-c")
			w := httptest.NewRecorder()
			sessionCreationTestRouter(t, store).ServeHTTP(w, r)
			if w.Code != 200 || strings.TrimSpace(w.Body.String()) != `{"sessions":[]}` {
				t.Fatalf("empty list = %d %s, want 200 with empty JSON array", w.Code, w.Body.String())
			}
		})
	}
}

func TestListSessionsHandlerRejectsInvalidAuthentication(t *testing.T) {
	store := NewMemorySessionStore()
	want := syntheticSession(t, "synthetic-private-session")
	if err := store.Add(want); err != nil {
		t.Fatal(err)
	}
	router := sessionCreationTestRouter(t, store)
	for _, values := range [][]string{nil, {"Bearer synthetic-unknown"}, {"Bearer synthetic-token-a", "Bearer synthetic-token-a"}} {
		r := httptest.NewRequest("GET", "/v1/sessions", nil)
		for _, value := range values {
			r.Header.Add("Authorization", value)
		}
		w := httptest.NewRecorder()
		router.ServeHTTP(w, r)
		response := w.Result()
		response.Body.Close()
		if response.StatusCode != 401 || response.Header.Get("Content-Type") != "application/json" || response.Header.Get("WWW-Authenticate") != "Bearer" || strings.TrimSpace(w.Body.String()) != `{"code":"unauthorized"}` {
			t.Errorf("invalid auth response: %d %s", response.StatusCode, w.Body.String())
		}
	}
	if got, exists := store.Get(want.ID); !exists || got != want || sessionCountForTest(store) != 1 {
		t.Error("unauthenticated list changed storage")
	}
}

func TestCreatedSessionAppearsInOwnerList(t *testing.T) {
	store := NewMemorySessionStore()
	router := sessionCreationTestRouter(t, store)
	create := httptest.NewRequest("POST", "/v1/sessions", strings.NewReader(`{"schema_version":"1","policy_id":"full-history-v1"}`))
	create.Header.Set("Authorization", "Bearer synthetic-token-a")
	created := httptest.NewRecorder()
	router.ServeHTTP(created, create)
	if created.Code != 201 {
		t.Fatalf("create status = %d", created.Code)
	}
	var expected CreateSessionResponse
	if err := json.Unmarshal(created.Body.Bytes(), &expected); err != nil {
		t.Fatal(err)
	}
	r := httptest.NewRequest("GET", "/v1/sessions", nil)
	r.Header.Set("Authorization", "Bearer synthetic-token-a")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	var got ListSessionsResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if w.Code != 200 || len(got.Sessions) != 1 || got.Sessions[0] != expected {
		t.Fatalf("created session missing or changed in list: %+v", got)
	}
}
