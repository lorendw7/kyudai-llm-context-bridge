package main

import (
	"errors"
	"fmt"
)

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
