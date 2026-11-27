package main

import (
	"errors"
	"time"
)

// L4-4: atomic admission of a turn.
//
// This file is the skeleton prepared by the assistant. The three method
// bodies below are the user's to write; the tests in turn_admit_test.go
// describe the agreed behavior. The store's internal layout for messages,
// turns and the active turn is the user's design and belongs in
// session_store.go.

// ErrSessionNotFound reports that no session has the given ID.
var ErrSessionNotFound = errors.New("session not found")

// ErrSessionBusy reports that a different request is still active in the
// session. No queue is kept; the caller decides when to try again.
var ErrSessionBusy = errors.New("session busy")

// AdmitTurn records the request, its user message and an accepted turn as
// one atomic operation.
//
//   - New request in a free session: returns the new turn, created = true.
//   - Same request ID and same payload: returns the existing turn,
//     created = false, and changes nothing.
//   - Same request ID, different payload: an error matching
//     ErrTurnRequestConflict.
//   - Different request ID while a turn is active: an error matching
//     ErrSessionBusy.
//
// The user message takes messageID, the session owner's user ID and now.
// Its Sequence equals the session Revision reached by this append.
func (s *MemorySessionStore) AdmitTurn(
	sessionID string,
	messageID string,
	request CreateTurnRequest,
	now time.Time,
) (turn Turn, created bool, err error) {
	// TODO(L4-4, user): implement.
	return Turn{}, false, errors.New("AdmitTurn is not implemented yet (L4-4)")
}

// Messages returns the session's messages in conversation order. The
// result is a copy: changing it must not change the store. An unknown
// session yields an empty, non-nil slice.
func (s *MemorySessionStore) Messages(sessionID string) []Message {
	// TODO(L4-4, user): implement.
	return []Message{}
}

// Turns returns the session's turns in admission order, as a copy. An
// unknown session yields an empty, non-nil slice.
func (s *MemorySessionStore) Turns(sessionID string) []Turn {
	// TODO(L4-4, user): implement.
	return []Turn{}
}
