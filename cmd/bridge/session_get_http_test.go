package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetSessionHandlerEnforcesOwnerScope(t *testing.T) {
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
	sessions := make([]Session, len(owners))
	for i, owner := range owners {
		session := syntheticSession(t, "synthetic-session-"+owner.lab+"-"+owner.user)
		session.LabID, session.OwnerUserID = owner.lab, owner.user
		if err := store.Add(session); err != nil {
			t.Fatal(err)
		}
		sessions[i] = session
	}
	router := sessionCreationTestRouter(t, store)
	for _, caller := range owners {
		for _, target := range sessions {
			t.Run(caller.lab+"/"+caller.user+"/"+target.ID, func(t *testing.T) {
				// Query/header assertions cannot override the token-derived identity.
				r := httptest.NewRequest("GET", "/v1/sessions/"+target.ID+"?user_id="+target.OwnerUserID+"&lab_id="+target.LabID, nil)
				r.Header.Set("Authorization", "Bearer "+caller.token)
				r.Header.Set("X-User-ID", target.OwnerUserID)
				r.Header.Set("X-Lab-ID", target.LabID)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, r)
				response := w.Result()
				defer response.Body.Close()
				if response.Header.Get("Content-Type") != "application/json" || response.Header.Get("Location") != "" {
					t.Errorf("unexpected committed headers: %v", response.Header)
				}
				if caller.lab != target.LabID || caller.user != target.OwnerUserID {
					if response.StatusCode != http.StatusNotFound || strings.TrimSpace(w.Body.String()) != `{"code":"session_not_found"}` {
						t.Fatalf("foreign session exposed: %d %s", response.StatusCode, w.Body.String())
					}
					missing := httptest.NewRequest("GET", "/v1/sessions/synthetic-absent", nil)
					missing.Header.Set("Authorization", "Bearer "+caller.token)
					missingResponse := httptest.NewRecorder()
					router.ServeHTTP(missingResponse, missing)
					if missingResponse.Code != response.StatusCode || missingResponse.Body.String() != w.Body.String() {
						t.Error("absent and foreign sessions have different status/body")
					}
					return
				}
				if response.StatusCode != http.StatusOK {
					t.Fatalf("owner GET status = %d, want 200", response.StatusCode)
				}
				var got CreateSessionResponse
				if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
					t.Fatal(err)
				}
				if got.ID != target.ID || got.PolicyID != target.PolicyID || got.Revision != target.Revision || !got.CreatedAt.Equal(target.CreatedAt) {
					t.Errorf("GET differs from saved session: %+v / %+v", got, target)
				}
				var fields map[string]json.RawMessage
				if err := json.Unmarshal(w.Body.Bytes(), &fields); err != nil || len(fields) != 4 {
					t.Error("GET exposed unexpected fields")
				}
			})
		}
	}
	if sessionCountForTest(store) != len(sessions) {
		t.Error("GET requests changed session count")
	}
	for _, want := range sessions {
		if got, exists := store.Get(want.ID); !exists || got != want {
			t.Errorf("GET requests changed stored session %s", want.ID)
		}
	}
}

func TestGetSessionHandlerRequiresAuthenticationForExistingAndMissingIDs(t *testing.T) {
	store := NewMemorySessionStore()
	session := syntheticSession(t, "synthetic-existing")
	if err := store.Add(session); err != nil {
		t.Fatal(err)
	}
	router := sessionCreationTestRouter(t, store)
	for _, id := range []string{session.ID, "synthetic-absent"} {
		for _, auth := range []string{"", "Bearer synthetic-unknown"} {
			r := httptest.NewRequest("GET", "/v1/sessions/"+id, nil)
			if auth != "" {
				r.Header.Set("Authorization", auth)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, r)
			response := w.Result()
			response.Body.Close()
			if response.StatusCode != http.StatusUnauthorized || strings.TrimSpace(w.Body.String()) != `{"code":"unauthorized"}` || response.Header.Get("WWW-Authenticate") != "Bearer" {
				t.Errorf("unauthenticated GET returned %d %s", response.StatusCode, w.Body.String())
			}
		}
	}
	if got, exists := store.Get(session.ID); !exists || got != session || sessionCountForTest(store) != 1 {
		t.Error("unauthenticated GET changed storage")
	}
}

func TestCreatedSessionCanBeReadAtLocation(t *testing.T) {
	store := NewMemorySessionStore()
	router := sessionCreationTestRouter(t, store)
	create := httptest.NewRequest("POST", "/v1/sessions", strings.NewReader(`{"schema_version":"1","policy_id":"full-history-v1"}`))
	create.Header.Set("Authorization", "Bearer synthetic-token-a")
	created := httptest.NewRecorder()
	router.ServeHTTP(created, create)
	response := created.Result()
	defer response.Body.Close()
	location := response.Header.Get("Location")
	if response.StatusCode != 201 || !strings.HasPrefix(location, "/v1/sessions/") {
		t.Fatalf("invalid creation response: %d %q", response.StatusCode, location)
	}
	get := httptest.NewRequest("GET", location, nil)
	get.Header.Set("Authorization", "Bearer synthetic-token-a")
	retrieved := httptest.NewRecorder()
	router.ServeHTTP(retrieved, get)
	if retrieved.Code != 200 || retrieved.Body.String() != created.Body.String() {
		t.Fatalf("create/GET mismatch: %d %s", retrieved.Code, retrieved.Body.String())
	}
}
