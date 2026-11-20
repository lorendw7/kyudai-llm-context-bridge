package main

import (
	"errors"
	"testing"
	"time"
)

const syntheticTurnSessionID = "synthetic-session-001"

var syntheticTurnCreatedAt = time.Date(2026, 11, 20, 9, 0, 0, 0, time.UTC)

// The correction sent to model B in the shared synthetic conversation.
func syntheticCorrectionRequest() CreateTurnRequest {
	return CreateTurnRequest{
		SchemaVersion: "1",
		RequestID:     "synthetic-turn-002",
		ModelID:       "local-model-b",
		Content:       "Correction: the deadline is November 15.",
	}
}

func newSyntheticTurn(t *testing.T) Turn {
	t.Helper()

	turn, err := NewTurn(syntheticTurnSessionID, syntheticCorrectionRequest(),
		syntheticTurnCreatedAt)
	if err != nil {
		t.Fatalf("NewTurn with a valid request: %v", err)
	}

	return turn
}

func TestNewTurnCopiesRequestAndStartsAccepted(t *testing.T) {
	request := syntheticCorrectionRequest()
	turn := newSyntheticTurn(t)

	want := Turn{
		SessionID:     syntheticTurnSessionID,
		RequestID:     request.RequestID,
		SchemaVersion: request.SchemaVersion,
		ModelID:       request.ModelID,
		Content:       request.Content,
		State:         TurnStateAccepted,
		CreatedAt:     syntheticTurnCreatedAt,
	}

	if turn != want {
		t.Fatalf("NewTurn = %+v\nwant      %+v", turn, want)
	}
}

func TestNewTurnKeepsContentExactly(t *testing.T) {
	request := syntheticCorrectionRequest()
	request.Content = "  Correction:\nthe deadline is November 15.  "

	turn, err := NewTurn(syntheticTurnSessionID, request, syntheticTurnCreatedAt)
	if err != nil {
		t.Fatalf("NewTurn: %v", err)
	}

	if turn.Content != request.Content {
		t.Fatalf("Content = %q, want the submitted text %q unchanged",
			turn.Content, request.Content)
	}
}

func TestNewTurnRejectsInvalidInput(t *testing.T) {
	tests := []struct {
		name      string
		sessionID string
		change    func(*CreateTurnRequest)
		createdAt time.Time
	}{
		{"empty_session_id", "", nil, syntheticTurnCreatedAt},
		{"blank_session_id", "   ", nil, syntheticTurnCreatedAt},
		{"zero_creation_time", syntheticTurnSessionID, nil, time.Time{}},
		{"unsupported_schema", syntheticTurnSessionID,
			func(r *CreateTurnRequest) { r.SchemaVersion = "2" }, syntheticTurnCreatedAt},
		{"empty_request_id", syntheticTurnSessionID,
			func(r *CreateTurnRequest) { r.RequestID = "" }, syntheticTurnCreatedAt},
		{"empty_model_id", syntheticTurnSessionID,
			func(r *CreateTurnRequest) { r.ModelID = " " }, syntheticTurnCreatedAt},
		{"empty_content", syntheticTurnSessionID,
			func(r *CreateTurnRequest) { r.Content = "\n" }, syntheticTurnCreatedAt},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := syntheticCorrectionRequest()
			if test.change != nil {
				test.change(&request)
			}

			turn, err := NewTurn(test.sessionID, request, test.createdAt)
			if err == nil {
				t.Fatalf("NewTurn accepted invalid input and returned %+v", turn)
			}

			if turn != (Turn{}) {
				t.Fatalf("NewTurn returned %+v with an error, want the zero Turn", turn)
			}
		})
	}
}

func TestCheckReplayAcceptsIdenticalRequest(t *testing.T) {
	turn := newSyntheticTurn(t)

	if err := turn.CheckReplay(syntheticCorrectionRequest()); err != nil {
		t.Fatalf("identical replay rejected: %v", err)
	}
}

// Option A from the 2026-11-20 decision: any byte difference in the compared
// fields is a change, including whitespace and case.
func TestCheckReplayRejectsChangedRequestAsConflict(t *testing.T) {
	tests := []struct {
		name   string
		change func(*CreateTurnRequest)
	}{
		{"model_switched_to_a", func(r *CreateTurnRequest) { r.ModelID = "local-model-a" }},
		{"model_case_changed", func(r *CreateTurnRequest) { r.ModelID = "Local-Model-B" }},
		{"content_changed", func(r *CreateTurnRequest) {
			r.Content = "Correction: the deadline is November 16."
		}},
		{"content_leading_space", func(r *CreateTurnRequest) { r.Content = " " + r.Content }},
		{"content_trailing_newline", func(r *CreateTurnRequest) { r.Content += "\n" }},
		{"content_case_changed", func(r *CreateTurnRequest) {
			r.Content = "correction: the deadline is November 15."
		}},
		{"model_and_content_changed", func(r *CreateTurnRequest) {
			r.ModelID = "local-model-a"
			r.Content = "When is the deadline?"
		}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			turn := newSyntheticTurn(t)
			request := syntheticCorrectionRequest()
			test.change(&request)

			err := turn.CheckReplay(request)
			if err == nil {
				t.Fatal("changed request accepted as a replay")
			}

			if !errors.Is(err, ErrTurnRequestConflict) {
				t.Fatalf("error %v does not match ErrTurnRequestConflict", err)
			}
		})
	}
}

// A request with another ID is not a replay of this turn, and an invalid
// request is invalid before it can be a conflict. Neither may be reported as
// a conflict.
func TestCheckReplayRejectsNonReplaysWithoutConflict(t *testing.T) {
	tests := []struct {
		name   string
		change func(*CreateTurnRequest)
	}{
		{"different_request_id", func(r *CreateTurnRequest) { r.RequestID = "synthetic-turn-003" }},
		{"request_id_case_changed", func(r *CreateTurnRequest) { r.RequestID = "Synthetic-Turn-002" }},
		{"unsupported_schema", func(r *CreateTurnRequest) { r.SchemaVersion = "2" }},
		{"empty_model_id", func(r *CreateTurnRequest) { r.ModelID = "" }},
		{"empty_content", func(r *CreateTurnRequest) { r.Content = "" }},
		{"empty_request_id", func(r *CreateTurnRequest) { r.RequestID = "" }},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			turn := newSyntheticTurn(t)
			request := syntheticCorrectionRequest()
			test.change(&request)

			err := turn.CheckReplay(request)
			if err == nil {
				t.Fatal("non-replay accepted")
			}

			if errors.Is(err, ErrTurnRequestConflict) {
				t.Fatalf("non-replay reported as a conflict: %v", err)
			}
		})
	}
}

// Replay is judged on the request alone. A failed or cancelled turn still
// answers an identical replay, so the client learns the original outcome.
func TestCheckReplayIgnoresTurnState(t *testing.T) {
	paths := map[string][]TurnState{
		"generating": {TurnStateGenerating},
		"complete":   {TurnStateGenerating, TurnStateAnswerSaved, TurnStateComplete},
		"failed":     {TurnStateGenerating, TurnStateFailed},
		"cancelled":  {TurnStateCancelled},
	}

	for name, path := range paths {
		t.Run(name, func(t *testing.T) {
			turn := newSyntheticTurn(t)

			for _, next := range path {
				var err error
				turn, err = turn.WithState(next)
				if err != nil {
					t.Fatalf("WithState(%s): %v", next, err)
				}
			}

			if err := turn.CheckReplay(syntheticCorrectionRequest()); err != nil {
				t.Fatalf("identical replay of a %s turn rejected: %v", name, err)
			}

			changed := syntheticCorrectionRequest()
			changed.ModelID = "local-model-a"

			if err := turn.CheckReplay(changed); !errors.Is(err, ErrTurnRequestConflict) {
				t.Fatalf("changed replay of a %s turn = %v, want a conflict", name, err)
			}
		})
	}
}

func TestWithStateReturnsNewTurnAndKeepsOriginal(t *testing.T) {
	original := newSyntheticTurn(t)
	snapshot := original

	updated, err := original.WithState(TurnStateGenerating)
	if err != nil {
		t.Fatalf("accepted -> generating: %v", err)
	}

	if original != snapshot {
		t.Fatalf("WithState changed the original turn to %+v", original)
	}

	want := snapshot
	want.State = TurnStateGenerating

	if updated != want {
		t.Fatalf("updated = %+v\nwant      %+v", updated, want)
	}
}

func TestWithStateRejectsDisallowedTransition(t *testing.T) {
	original := newSyntheticTurn(t)
	snapshot := original

	// The lesson's opening question: a reply cannot be cancelled once saved,
	// and accepted cannot jump straight to answer_saved.
	updated, err := original.WithState(TurnStateAnswerSaved)
	if !errors.Is(err, ErrInvalidTurnTransition) {
		t.Fatalf("accepted -> answer_saved = %v, want ErrInvalidTurnTransition", err)
	}

	if updated != (Turn{}) {
		t.Fatalf("WithState returned %+v with an error, want the zero Turn", updated)
	}

	if original != snapshot {
		t.Fatalf("failed WithState changed the original turn to %+v", original)
	}
}

func TestWithStateRejectsInvalidStateValue(t *testing.T) {
	turn := newSyntheticTurn(t)

	updated, err := turn.WithState("running")
	if err == nil {
		t.Fatalf("invalid state accepted, returned %+v", updated)
	}

	if errors.Is(err, ErrInvalidTurnTransition) {
		t.Fatalf("invalid value reported as a transition error: %v", err)
	}

	if updated != (Turn{}) {
		t.Fatalf("WithState returned %+v with an error, want the zero Turn", updated)
	}
}

func TestWithStateFollowsFullHistoryPath(t *testing.T) {
	turn := newSyntheticTurn(t)
	path := []TurnState{TurnStateGenerating, TurnStateAnswerSaved, TurnStateComplete}

	for _, next := range path {
		previous := turn

		var err error
		turn, err = turn.WithState(next)
		if err != nil {
			t.Fatalf("%s -> %s: %v", previous.State, next, err)
		}

		if turn.State != next {
			t.Fatalf("State = %s, want %s", turn.State, next)
		}

		// Only the state changes along the path.
		previous.State = next
		if turn != previous {
			t.Fatalf("WithState(%s) changed other fields: %+v", next, turn)
		}
	}

	if _, err := turn.WithState(TurnStateCancelled); !errors.Is(err, ErrInvalidTurnTransition) {
		t.Fatalf("complete -> cancelled = %v, want ErrInvalidTurnTransition", err)
	}
}

// The zero Turn is not a valid turn. Construction goes through NewTurn.
func TestZeroTurnCannotChangeState(t *testing.T) {
	var turn Turn

	if err := turn.State.Validate(); err == nil {
		t.Fatal("the zero Turn has a valid state")
	}

	if _, err := turn.WithState(TurnStateGenerating); err == nil {
		t.Fatal("the zero Turn changed state")
	}
}
