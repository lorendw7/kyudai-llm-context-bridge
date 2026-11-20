package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var ErrTurnRequestConflict = errors.New("turn request conflict")

type Turn struct {
	SessionID     string
	RequestID     string
	SchemaVersion string
	ModelID       string
	Content       string
	State         TurnState
	CreatedAt     time.Time
}

func NewTurn(
	sessionID string,
	request CreateTurnRequest,
	createdAt time.Time,
) (Turn, error) {
	if strings.TrimSpace(sessionID) == "" {
		return Turn{}, errors.New("session ID is required")
	}

	if err := request.Validate(); err != nil {
		return Turn{}, err
	}

	if createdAt.IsZero() {
		return Turn{}, errors.New("creation time is required")
	}

	return Turn{
		SessionID:     sessionID,
		RequestID:     request.RequestID,
		SchemaVersion: request.SchemaVersion,
		ModelID:       request.ModelID,
		Content:       request.Content,
		State:         TurnStateAccepted,
		CreatedAt:     createdAt,
	}, nil
}

func (t Turn) CheckReplay(request CreateTurnRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}

	if request.RequestID != t.RequestID {
		return fmt.Errorf("request %q is not a replay of %q",
			request.RequestID, t.RequestID)
	}

	if request.SchemaVersion != t.SchemaVersion ||
		request.ModelID != t.ModelID ||
		request.Content != t.Content {
		return fmt.Errorf("%w: %s", ErrTurnRequestConflict, t.RequestID)
	}

	return nil
}

func (t Turn) WithState(next TurnState) (Turn, error) {
	if err := t.State.TransitionTo(next); err != nil {
		return Turn{}, err
	}

	t.State = next
	return t, nil
}

var ErrInvalidTurnTransition = errors.New("invalid turn transition")

type TurnState string

const (
	TurnStateAccepted         TurnState = "accepted"
	TurnStateGenerating       TurnState = "generating"
	TurnStateAnswerSaved      TurnState = "answer_saved"
	TurnStateMemoryPending    TurnState = "memory_pending"
	TurnStateComplete         TurnState = "complete"
	TurnStateFailed           TurnState = "failed"
	TurnStateCancelled        TurnState = "cancelled"
	TurnStateRecoveryRequired TurnState = "recovery_required"
)

func (state TurnState) Validate() error {
	switch state {
	case TurnStateAccepted,
		TurnStateGenerating,
		TurnStateAnswerSaved,
		TurnStateMemoryPending,
		TurnStateComplete,
		TurnStateFailed,
		TurnStateCancelled,
		TurnStateRecoveryRequired:
		return nil
	default:
		return fmt.Errorf("invalid turn state: %q", state)
	}
}

func (state TurnState) IsTerminal() bool {
	switch state {
	case TurnStateComplete, TurnStateFailed, TurnStateCancelled:
		return true
	default:
		return false
	}
}

var turnTransitions = map[TurnState][]TurnState{
	TurnStateAccepted:         {TurnStateGenerating, TurnStateCancelled},
	TurnStateGenerating:       {TurnStateAnswerSaved, TurnStateFailed, TurnStateCancelled, TurnStateRecoveryRequired},
	TurnStateAnswerSaved:      {TurnStateMemoryPending, TurnStateComplete},
	TurnStateMemoryPending:    {TurnStateComplete, TurnStateRecoveryRequired},
	TurnStateRecoveryRequired: {TurnStateMemoryPending, TurnStateFailed},
}

func (state TurnState) TransitionTo(next TurnState) error {
	if err := state.Validate(); err != nil {
		return err
	}

	if err := next.Validate(); err != nil {
		return err
	}

	for _, allowed := range turnTransitions[state] {
		if allowed == next {
			return nil
		}
	}

	return fmt.Errorf("%w: %s -> %s", ErrInvalidTurnTransition, state, next)
}
