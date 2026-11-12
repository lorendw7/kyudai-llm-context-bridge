package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func syntheticSession(t *testing.T, id string) Session {
	t.Helper()
	s, err := NewSession(id, "lab-a", "user-1", "full-history-v1", time.Date(2026, 11, 12, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestMemorySessionStoreRoundTripAndValueIsolation(t *testing.T) {
	store := NewMemorySessionStore()
	want := syntheticSession(t, "synthetic-session-1")
	input := want
	if err := store.Add(input); err != nil {
		t.Fatal(err)
	}
	input.OwnerUserID = "changed-input"
	got, exists := store.Get(want.ID)
	if !exists || got != want {
		t.Fatalf("Get = %+v, %v; want %+v, true", got, exists, want)
	}
	got.OwnerUserID = "changed-output"
	if again, exists := store.Get(want.ID); !exists || again != want {
		t.Fatal("modifying a returned Session changed stored ownership")
	}
	if missing, exists := store.Get("missing-session"); exists || missing != (Session{}) {
		t.Fatalf("unknown session returned %+v, %v", missing, exists)
	}
	if _, exists := NewMemorySessionStore().Get(want.ID); exists {
		t.Fatal("separate stores shared session data")
	}
}

func TestMemorySessionStoreDuplicateDoesNotOverwriteOrKeepLock(t *testing.T) {
	store := NewMemorySessionStore()
	want := syntheticSession(t, "synthetic-session-1")
	if err := store.Add(want); err != nil {
		t.Fatal(err)
	}
	duplicate := want
	duplicate.OwnerUserID = "user-2"
	if err := store.Add(duplicate); err == nil {
		t.Fatal("duplicate ID was accepted")
	}
	if got, exists := store.Get(want.ID); !exists || got != want {
		t.Fatal("duplicate overwrote the original session")
	}
	if err := store.Add(syntheticSession(t, "synthetic-session-2")); err != nil {
		t.Fatal(err)
	}
}

func TestMemorySessionStoreConcurrentDuplicateHasOneWinner(t *testing.T) {
	store := NewMemorySessionStore()
	want := syntheticSession(t, "synthetic-shared-session")
	const workers = 32
	start := make(chan struct{})
	results := make(chan error, workers)
	for i := 0; i < workers; i++ {
		go func() {
			<-start
			results <- store.Add(want)
		}()
	}
	close(start)
	successes := 0
	for i := 0; i < workers; i++ {
		if err := <-results; err == nil {
			successes++
		}
	}
	if successes != 1 {
		t.Fatalf("successful concurrent inserts = %d, want 1", successes)
	}
	if got, exists := store.Get(want.ID); !exists || got != want {
		t.Fatal("winning session was not preserved")
	}
}

func TestMemorySessionStoreConcurrentReadsAndWrites(t *testing.T) {
	store := NewMemorySessionStore()
	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		session := syntheticSession(t, fmt.Sprintf("synthetic-session-%d", i))
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := store.Add(session); err != nil {
				t.Error(err)
				return
			}
			if got, exists := store.Get(session.ID); !exists || got != session {
				t.Errorf("concurrent round trip failed for %s", session.ID)
			}
		}()
	}
	wg.Wait()
	for i := 0; i < 32; i++ {
		if _, exists := store.Get(fmt.Sprintf("synthetic-session-%d", i)); !exists {
			t.Errorf("session %d was lost", i)
		}
	}
}

func TestMemorySessionStoreListOwnedFiltersBothScopeFields(t *testing.T) {
	store := NewMemorySessionStore()
	owners := []Principal{
		{LabID: "lab-a", UserID: "user-1"},
		{LabID: "lab-a", UserID: "user-2"},
		{LabID: "lab-b", UserID: "user-1"},
		{LabID: "lab-b", UserID: "user-2"},
	}
	want := make(map[Principal]map[string]Session)
	for i, owner := range owners {
		want[owner] = make(map[string]Session)
		for j := 0; j < 2; j++ {
			session := syntheticSession(t, fmt.Sprintf("synthetic-owned-%d-%d", i, j))
			session.LabID, session.OwnerUserID = owner.LabID, owner.UserID
			if err := store.Add(session); err != nil {
				t.Fatal(err)
			}
			want[owner][session.ID] = session
		}
	}
	for _, owner := range owners {
		got := store.ListOwned(owner)
		if len(got) != len(want[owner]) {
			t.Fatalf("list for %+v has %d sessions, want %d", owner, len(got), len(want[owner]))
		}
		seen := make(map[string]bool)
		for _, session := range got {
			expected, exists := want[owner][session.ID]
			if !exists || session != expected || seen[session.ID] {
				t.Errorf("list for %+v returned foreign, altered or repeated session: %+v", owner, session)
			}
			seen[session.ID] = true
		}
	}
	for _, emptyStore := range []*MemorySessionStore{store, NewMemorySessionStore()} {
		got := emptyStore.ListOwned(Principal{LabID: "lab-c", UserID: "user-3"})
		if got == nil || len(got) != 0 {
			t.Fatalf("empty result = %#v, want initialized empty slice", got)
		}
	}
}

func TestMemorySessionStoreListOwnedReturnsIndependentValues(t *testing.T) {
	store := NewMemorySessionStore()
	want := syntheticSession(t, "synthetic-owned-session")
	owner := Principal{LabID: want.LabID, UserID: want.OwnerUserID}
	if err := store.Add(want); err != nil {
		t.Fatal(err)
	}
	first := store.ListOwned(owner)
	second := store.ListOwned(owner)
	if len(first) != 1 || len(second) != 1 {
		t.Fatal("expected one owned session")
	}
	first[0].OwnerUserID = "synthetic-changed-owner"
	first[0].PolicyID = "synthetic-changed-policy"
	if second[0] != want {
		t.Error("separate list results share mutable elements")
	}
	if stored, exists := store.Get(want.ID); !exists || stored != want {
		t.Error("changing a list result changed stored data")
	}
	if again := store.ListOwned(owner); len(again) != 1 || again[0] != want {
		t.Error("changing a list result affected subsequent listing")
	}
}

func TestMemorySessionStoreListOwnedDuringConcurrentWrites(t *testing.T) {
	store := NewMemorySessionStore()
	owner := Principal{LabID: "lab-a", UserID: "user-1"}
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 0; i < 32; i++ {
		session := syntheticSession(t, fmt.Sprintf("synthetic-list-concurrent-%d", i))
		if i%2 == 1 {
			session.OwnerUserID = "user-2"
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			if err := store.Add(session); err != nil {
				t.Error(err)
			}
			for _, listed := range store.ListOwned(owner) {
				if !listed.OwnedBy(owner.LabID, owner.UserID) {
					t.Errorf("concurrent list exposed another owner: %+v", listed)
				}
			}
		}()
	}
	close(start)
	wg.Wait()
	if got := len(store.ListOwned(owner)); got != 16 {
		t.Errorf("final owned count = %d, want 16", got)
	}
}
