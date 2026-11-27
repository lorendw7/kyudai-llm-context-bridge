package main

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

var syntheticAdmitTime = time.Date(2026, 11, 27, 10, 0, 0, 0, time.UTC)

// The first instruction of the shared synthetic conversation, sent to A.
func syntheticDeadlineRequest() CreateTurnRequest {
	return CreateTurnRequest{
		SchemaVersion: "1",
		RequestID:     "synthetic-turn-001",
		ModelID:       "local-model-a",
		Content:       "The deadline is November 12.",
	}
}

func newAdmitStore(t *testing.T, sessionIDs ...string) *MemorySessionStore {
	t.Helper()

	store := NewMemorySessionStore()
	for _, id := range sessionIDs {
		if err := store.Add(syntheticSession(t, id)); err != nil {
			t.Fatalf("Add(%q): %v", id, err)
		}
	}

	return store
}

type admitSnapshot struct {
	revision uint64
	messages []Message
	turns    []Turn
}

func snapshotSession(t *testing.T, store *MemorySessionStore, sessionID string) admitSnapshot {
	t.Helper()

	session, ok := store.Get(sessionID)
	if !ok {
		t.Fatalf("Get(%q): session is missing", sessionID)
	}

	return admitSnapshot{
		revision: session.Revision,
		messages: store.Messages(sessionID),
		turns:    store.Turns(sessionID),
	}
}

func requireUnchanged(t *testing.T, store *MemorySessionStore, sessionID string, before admitSnapshot) {
	t.Helper()

	after := snapshotSession(t, store, sessionID)
	if after.revision != before.revision {
		t.Errorf("Revision changed from %d to %d", before.revision, after.revision)
	}
	if len(after.messages) != len(before.messages) {
		t.Errorf("message count changed from %d to %d", len(before.messages), len(after.messages))
	}
	if len(after.turns) != len(before.turns) {
		t.Errorf("turn count changed from %d to %d", len(before.turns), len(after.turns))
	}
	for i := range before.messages {
		if i < len(after.messages) && after.messages[i] != before.messages[i] {
			t.Errorf("message %d changed:\n got %+v\nwant %+v", i, after.messages[i], before.messages[i])
		}
	}
	for i := range before.turns {
		if i < len(after.turns) && after.turns[i] != before.turns[i] {
			t.Errorf("turn %d changed:\n got %+v\nwant %+v", i, after.turns[i], before.turns[i])
		}
	}
}

func TestAdmitTurnRecordsRequestMessageAndAcceptedTogether(t *testing.T) {
	const sessionID = "synthetic-session-001"
	store := newAdmitStore(t, sessionID)
	session, _ := store.Get(sessionID)
	request := syntheticDeadlineRequest()

	if session.Revision != 1 {
		t.Fatalf("a new session starts at Revision %d, want 1", session.Revision)
	}
	if got := store.Messages(sessionID); got == nil || len(got) != 0 {
		t.Fatalf("Messages before admission = %#v, want empty non-nil", got)
	}
	if got := store.Turns(sessionID); got == nil || len(got) != 0 {
		t.Fatalf("Turns before admission = %#v, want empty non-nil", got)
	}

	turn, created, err := store.AdmitTurn(sessionID, "synthetic-message-user1", request, syntheticAdmitTime)
	if err != nil {
		t.Fatalf("AdmitTurn: %v", err)
	}
	if !created {
		t.Fatal("created = false for a new request, want true")
	}

	wantTurn := Turn{
		SessionID:     sessionID,
		RequestID:     request.RequestID,
		SchemaVersion: request.SchemaVersion,
		ModelID:       request.ModelID,
		Content:       request.Content,
		State:         TurnStateAccepted,
		CreatedAt:     syntheticAdmitTime,
	}
	if turn != wantTurn {
		t.Fatalf("turn = %+v\nwant   %+v", turn, wantTurn)
	}

	wantMessage := Message{
		ID:        "synthetic-message-user1",
		SessionID: sessionID,
		UserID:    session.OwnerUserID,
		Role:      MessageRoleUser,
		Content:   request.Content,
		CreatedAt: syntheticAdmitTime,
		Sequence:  2,
	}
	messages := store.Messages(sessionID)
	if len(messages) != 1 || messages[0] != wantMessage {
		t.Fatalf("Messages = %+v\nwant one message %+v", messages, wantMessage)
	}

	turns := store.Turns(sessionID)
	if len(turns) != 1 || turns[0] != wantTurn {
		t.Fatalf("Turns = %+v\nwant one turn %+v", turns, wantTurn)
	}

	after, _ := store.Get(sessionID)
	if after.Revision != 2 {
		t.Fatalf("Revision after one appended message = %d, want 2", after.Revision)
	}
	if messages[0].Sequence != after.Revision {
		t.Fatalf("message Sequence %d differs from session Revision %d",
			messages[0].Sequence, after.Revision)
	}
}

func TestAdmitTurnReplayReturnsSameTurnAndChangesNothing(t *testing.T) {
	const sessionID = "synthetic-session-001"
	store := newAdmitStore(t, sessionID)
	request := syntheticDeadlineRequest()

	first, _, err := store.AdmitTurn(sessionID, "synthetic-message-user1", request, syntheticAdmitTime)
	if err != nil {
		t.Fatalf("first AdmitTurn: %v", err)
	}
	before := snapshotSession(t, store, sessionID)

	// A replay arrives later and carries a different message ID. Neither
	// may replace what the first admission recorded.
	replay, created, err := store.AdmitTurn(sessionID, "synthetic-message-ignored",
		request, syntheticAdmitTime.Add(time.Minute))
	if err != nil {
		t.Fatalf("replay AdmitTurn: %v", err)
	}
	if created {
		t.Fatal("created = true for a replay, want false")
	}
	if replay != first {
		t.Fatalf("replay turn = %+v\nwant first  %+v", replay, first)
	}
	requireUnchanged(t, store, sessionID, before)
}

func TestAdmitTurnChangedPayloadConflicts(t *testing.T) {
	changes := map[string]func(*CreateTurnRequest){
		"model": func(r *CreateTurnRequest) { r.ModelID = "local-model-b" },
		"text":  func(r *CreateTurnRequest) { r.Content = "The deadline is November 15." },
		"space": func(r *CreateTurnRequest) { r.Content += " " },
	}

	for name, change := range changes {
		t.Run(name, func(t *testing.T) {
			const sessionID = "synthetic-session-001"
			store := newAdmitStore(t, sessionID)
			request := syntheticDeadlineRequest()

			if _, _, err := store.AdmitTurn(sessionID, "synthetic-message-user1",
				request, syntheticAdmitTime); err != nil {
				t.Fatalf("first AdmitTurn: %v", err)
			}
			before := snapshotSession(t, store, sessionID)

			changed := request
			change(&changed)
			turn, created, err := store.AdmitTurn(sessionID, "synthetic-message-other",
				changed, syntheticAdmitTime)

			if !errors.Is(err, ErrTurnRequestConflict) {
				t.Fatalf("err = %v, want one matching ErrTurnRequestConflict", err)
			}
			if errors.Is(err, ErrSessionBusy) {
				t.Fatal("a changed payload under the same request ID must be a conflict, not busy")
			}
			if created || turn != (Turn{}) {
				t.Fatalf("turn = %+v, created = %v; want zero values on conflict", turn, created)
			}
			requireUnchanged(t, store, sessionID, before)
		})
	}
}

func TestAdmitTurnDifferentRequestWhileActiveIsBusy(t *testing.T) {
	const sessionID = "synthetic-session-001"
	store := newAdmitStore(t, sessionID)

	if _, _, err := store.AdmitTurn(sessionID, "synthetic-message-user1",
		syntheticDeadlineRequest(), syntheticAdmitTime); err != nil {
		t.Fatalf("first AdmitTurn: %v", err)
	}
	before := snapshotSession(t, store, sessionID)

	turn, created, err := store.AdmitTurn(sessionID, "synthetic-message-user2",
		syntheticCorrectionRequest(), syntheticAdmitTime)

	if !errors.Is(err, ErrSessionBusy) {
		t.Fatalf("err = %v, want one matching ErrSessionBusy", err)
	}
	if errors.Is(err, ErrTurnRequestConflict) {
		t.Fatal("a different request ID must be busy, not a conflict")
	}
	if created || turn != (Turn{}) {
		t.Fatalf("turn = %+v, created = %v; want zero values when busy", turn, created)
	}
	requireUnchanged(t, store, sessionID, before)
}

func TestAdmitTurnRejectsUnknownSessionAndInvalidRequest(t *testing.T) {
	const sessionID = "synthetic-session-001"
	store := newAdmitStore(t, sessionID)
	before := snapshotSession(t, store, sessionID)

	_, created, err := store.AdmitTurn("synthetic-session-missing", "synthetic-message-user1",
		syntheticDeadlineRequest(), syntheticAdmitTime)
	if !errors.Is(err, ErrSessionNotFound) || created {
		t.Fatalf("unknown session: err = %v, created = %v; want ErrSessionNotFound", err, created)
	}
	if got := store.Messages("synthetic-session-missing"); got == nil || len(got) != 0 {
		t.Fatalf("Messages for an unknown session = %#v, want empty non-nil", got)
	}
	if got := store.Turns("synthetic-session-missing"); got == nil || len(got) != 0 {
		t.Fatalf("Turns for an unknown session = %#v, want empty non-nil", got)
	}

	invalid := map[string]CreateTurnRequest{}
	for name, change := range map[string]func(*CreateTurnRequest){
		"schema":     func(r *CreateTurnRequest) { r.SchemaVersion = "2" },
		"request ID": func(r *CreateTurnRequest) { r.RequestID = " " },
		"model":      func(r *CreateTurnRequest) { r.ModelID = "" },
		"content":    func(r *CreateTurnRequest) { r.Content = "\n" },
	} {
		request := syntheticDeadlineRequest()
		change(&request)
		invalid[name] = request
	}

	for name, request := range invalid {
		t.Run(name, func(t *testing.T) {
			turn, created, err := store.AdmitTurn(sessionID, "synthetic-message-user1",
				request, syntheticAdmitTime)
			if err == nil {
				t.Fatal("err = nil for an invalid request")
			}
			for _, sentinel := range []error{ErrSessionNotFound, ErrSessionBusy, ErrTurnRequestConflict} {
				if errors.Is(err, sentinel) {
					t.Fatalf("an invalid request reported %v", sentinel)
				}
			}
			if created || turn != (Turn{}) {
				t.Fatalf("turn = %+v, created = %v; want zero values", turn, created)
			}
		})
	}

	_, _, err = store.AdmitTurn(sessionID, "synthetic-message-user1",
		syntheticDeadlineRequest(), time.Time{})
	if err == nil {
		t.Fatal("err = nil for a zero admission time")
	}

	requireUnchanged(t, store, sessionID, before)
}

func TestAdmitTurnKeepsSessionsIndependent(t *testing.T) {
	const first, second = "synthetic-session-001", "synthetic-session-002"
	store := newAdmitStore(t, first, second)
	request := syntheticDeadlineRequest()

	if _, _, err := store.AdmitTurn(first, "synthetic-message-first", request, syntheticAdmitTime); err != nil {
		t.Fatalf("AdmitTurn in the first session: %v", err)
	}

	// The same request ID in another session is a different request, and an
	// active turn elsewhere does not make this session busy.
	turn, created, err := store.AdmitTurn(second, "synthetic-message-second", request, syntheticAdmitTime)
	if err != nil || !created {
		t.Fatalf("AdmitTurn in the second session: created = %v, err = %v", created, err)
	}
	if turn.SessionID != second {
		t.Fatalf("turn.SessionID = %q, want %q", turn.SessionID, second)
	}

	for _, id := range []string{first, second} {
		got := snapshotSession(t, store, id)
		if got.revision != 2 || len(got.messages) != 1 || len(got.turns) != 1 {
			t.Fatalf("%s: revision %d, %d messages, %d turns; want 2, 1, 1",
				id, got.revision, len(got.messages), len(got.turns))
		}
		if got.messages[0].SessionID != id || got.messages[0].Sequence != 2 {
			t.Fatalf("%s: message = %+v, want its own session and Sequence 2", id, got.messages[0])
		}
	}
}

func TestMessagesAndTurnsReturnCopies(t *testing.T) {
	const sessionID = "synthetic-session-001"
	store := newAdmitStore(t, sessionID)

	if _, _, err := store.AdmitTurn(sessionID, "synthetic-message-user1",
		syntheticDeadlineRequest(), syntheticAdmitTime); err != nil {
		t.Fatalf("AdmitTurn: %v", err)
	}
	before := snapshotSession(t, store, sessionID)

	messages := store.Messages(sessionID)
	messages[0].Content = "changed by the caller"
	messages = append(messages, Message{ID: "synthetic-message-forged"})
	_ = messages

	turns := store.Turns(sessionID)
	turns[0].State = TurnStateComplete
	turns = append(turns, Turn{RequestID: "synthetic-turn-forged"})
	_ = turns

	requireUnchanged(t, store, sessionID, before)

	// The forged terminal state must not have freed the session.
	_, _, err := store.AdmitTurn(sessionID, "synthetic-message-user2",
		syntheticCorrectionRequest(), syntheticAdmitTime)
	if !errors.Is(err, ErrSessionBusy) {
		t.Fatalf("err = %v after changing a returned turn, want ErrSessionBusy", err)
	}
}

// Thirty-two identical requests arrive at once: one turn, one message.
func TestAdmitTurnConcurrentIdenticalRequestsAdmitOnce(t *testing.T) {
	const sessionID = "synthetic-session-001"
	const callers = 32
	store := newAdmitStore(t, sessionID)
	request := syntheticDeadlineRequest()

	type result struct {
		turn    Turn
		created bool
		err     error
	}
	results := make([]result, callers)

	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < callers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			turn, created, err := store.AdmitTurn(sessionID,
				fmt.Sprintf("synthetic-message-%02d", i), request, syntheticAdmitTime)
			results[i] = result{turn, created, err}
		}(i)
	}
	close(start)
	wg.Wait()

	createdCount := 0
	for i, r := range results {
		if r.err != nil {
			t.Fatalf("caller %d: err = %v, want nil for an identical request", i, r.err)
		}
		if r.created {
			createdCount++
		}
		if r.turn != results[0].turn {
			t.Fatalf("caller %d received a different turn:\n got %+v\nwant %+v", i, r.turn, results[0].turn)
		}
	}
	if createdCount != 1 {
		t.Fatalf("%d callers saw created = true, want exactly 1", createdCount)
	}

	got := snapshotSession(t, store, sessionID)
	if len(got.turns) != 1 || len(got.messages) != 1 || got.revision != 2 {
		t.Fatalf("%d turns, %d messages, Revision %d; want 1, 1, 2",
			len(got.turns), len(got.messages), got.revision)
	}
}

// Two different requests compete, sixteen callers each. Which one wins is
// not fixed; the invariants are.
func TestAdmitTurnConcurrentCompetingRequestsHaveOneWinner(t *testing.T) {
	const sessionID = "synthetic-session-001"
	const callersPerRequest = 16
	store := newAdmitStore(t, sessionID)
	requests := []CreateTurnRequest{syntheticDeadlineRequest(), syntheticCorrectionRequest()}

	type result struct {
		created bool
		err     error
	}
	results := make([][]result, len(requests))

	start := make(chan struct{})
	var wg sync.WaitGroup
	for r := range requests {
		results[r] = make([]result, callersPerRequest)
		for i := 0; i < callersPerRequest; i++ {
			wg.Add(1)
			go func(r, i int) {
				defer wg.Done()
				<-start
				_, created, err := store.AdmitTurn(sessionID,
					fmt.Sprintf("synthetic-message-%d-%02d", r, i), requests[r], syntheticAdmitTime)
				results[r][i] = result{created, err}
			}(r, i)
		}
	}
	close(start)
	wg.Wait()

	winners := 0
	for r := range requests {
		admitted, created, busy := 0, 0, 0
		for i, res := range results[r] {
			switch {
			case res.err == nil:
				admitted++
				if res.created {
					created++
				}
			case errors.Is(res.err, ErrSessionBusy):
				busy++
			default:
				t.Fatalf("request %d caller %d: unexpected err %v", r, i, res.err)
			}
		}

		switch {
		case admitted == callersPerRequest && created == 1:
			winners++
		case busy == callersPerRequest:
		default:
			t.Fatalf("request %d: %d admitted (%d created), %d busy; want all admitted with one created, or all busy",
				r, admitted, created, busy)
		}
	}
	if winners != 1 {
		t.Fatalf("%d requests were admitted, want exactly 1", winners)
	}

	got := snapshotSession(t, store, sessionID)
	if len(got.turns) != 1 || len(got.messages) != 1 || got.revision != 2 {
		t.Fatalf("%d turns, %d messages, Revision %d; want 1, 1, 2",
			len(got.turns), len(got.messages), got.revision)
	}
	if got.messages[0].Content != got.turns[0].Content {
		t.Fatalf("the stored message %q does not belong to the stored turn %q",
			got.messages[0].Content, got.turns[0].Content)
	}
}

// Readers run while admissions are in flight, so the race detector can see
// unsynchronized access to the store's internals.
func TestAdmitTurnConcurrentReadersSeeConsistentState(t *testing.T) {
	const sessions = 8
	ids := make([]string, sessions)
	for i := range ids {
		ids[i] = fmt.Sprintf("synthetic-session-%03d", i+1)
	}
	store := newAdmitStore(t, ids...)

	start := make(chan struct{})
	var wg sync.WaitGroup
	for i, id := range ids {
		wg.Add(2)
		go func(i int, id string) {
			defer wg.Done()
			<-start
			if _, _, err := store.AdmitTurn(id, fmt.Sprintf("synthetic-message-%02d", i),
				syntheticDeadlineRequest(), syntheticAdmitTime); err != nil {
				t.Errorf("AdmitTurn(%s): %v", id, err)
			}
		}(i, id)
		go func(id string) {
			defer wg.Done()
			<-start
			for n := 0; n < 50; n++ {
				messages := store.Messages(id)
				turns := store.Turns(id)
				session, _ := store.Get(id)
				// Messages, Turns and Get are three separate reads, so they
				// may straddle the admission. Each alone must be 0 or 1.
				if len(messages) > 1 || len(turns) > 1 || session.Revision > 2 {
					t.Errorf("%s: %d messages, %d turns, Revision %d", id,
						len(messages), len(turns), session.Revision)
					return
				}
			}
		}(id)
	}
	close(start)
	wg.Wait()
}
