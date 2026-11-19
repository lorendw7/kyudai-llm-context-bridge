package main

import (
	"errors"
	"fmt"
	"testing"
)

// allowedTurnTransitions mirrors the state diagram in
// docs/contracts.md#turn-states. It is written independently of the
// implementation so the test fails if either drifts from the contract.
var allowedTurnTransitions = map[TurnState][]TurnState{
	TurnStateAccepted: {
		TurnStateGenerating,
		TurnStateCancelled,
	},
	TurnStateGenerating: {
		TurnStateAnswerSaved,
		TurnStateFailed,
		TurnStateCancelled,
		TurnStateRecoveryRequired,
	},
	TurnStateAnswerSaved: {
		TurnStateMemoryPending,
		TurnStateComplete,
	},
	TurnStateMemoryPending: {
		TurnStateComplete,
		TurnStateRecoveryRequired,
	},
	TurnStateRecoveryRequired: {
		TurnStateMemoryPending,
		TurnStateFailed,
	},
	TurnStateComplete:  {},
	TurnStateFailed:    {},
	TurnStateCancelled: {},
}

var allTurnStates = []TurnState{
	TurnStateAccepted,
	TurnStateGenerating,
	TurnStateAnswerSaved,
	TurnStateMemoryPending,
	TurnStateComplete,
	TurnStateFailed,
	TurnStateCancelled,
	TurnStateRecoveryRequired,
}

func isAllowedTurnTransition(from TurnState, to TurnState) bool {
	for _, next := range allowedTurnTransitions[from] {
		if next == to {
			return true
		}
	}

	return false
}

func TestTurnTransitionFixtureCoversEveryState(t *testing.T) {
	if len(allowedTurnTransitions) != len(allTurnStates) {
		t.Fatalf("fixture lists %d states, want %d",
			len(allowedTurnTransitions), len(allTurnStates))
	}

	for _, state := range allTurnStates {
		if err := state.Validate(); err != nil {
			t.Fatalf("fixture state %q is not canonical: %v", state, err)
		}

		if _, exists := allowedTurnTransitions[state]; !exists {
			t.Fatalf("fixture is missing state %q", state)
		}
	}
}

// Every one of the 64 ordered pairs, including a state to itself, must match
// the diagram exactly.
func TestTurnStateTransitionToMatchesContractDiagram(t *testing.T) {
	for _, from := range allTurnStates {
		for _, to := range allTurnStates {
			t.Run(fmt.Sprintf("%s_to_%s", from, to), func(t *testing.T) {
				err := from.TransitionTo(to)

				if isAllowedTurnTransition(from, to) {
					if err != nil {
						t.Fatalf("%s -> %s rejected: %v, want nil", from, to, err)
					}
					return
				}

				if err == nil {
					t.Fatalf("%s -> %s allowed, want an error", from, to)
				}

				if !errors.Is(err, ErrInvalidTurnTransition) {
					t.Fatalf("%s -> %s error %v does not match "+
						"ErrInvalidTurnTransition", from, to, err)
				}
			})
		}
	}
}

func TestTurnStateIsTerminal(t *testing.T) {
	terminal := map[TurnState]bool{
		TurnStateComplete:  true,
		TurnStateFailed:    true,
		TurnStateCancelled: true,
	}

	for _, state := range allTurnStates {
		t.Run(string(state), func(t *testing.T) {
			if got := state.IsTerminal(); got != terminal[state] {
				t.Fatalf("%s.IsTerminal() = %v, want %v",
					state, got, terminal[state])
			}
		})
	}

	for _, invalid := range []TurnState{"", "running", "Complete", " failed"} {
		t.Run("invalid_"+fmt.Sprintf("%q", string(invalid)), func(t *testing.T) {
			if invalid.IsTerminal() {
				t.Fatalf("invalid state %q reported as terminal", invalid)
			}
		})
	}
}

func TestTerminalTurnStatesHaveNoExit(t *testing.T) {
	for _, from := range allTurnStates {
		if !from.IsTerminal() {
			continue
		}

		for _, to := range allTurnStates {
			if err := from.TransitionTo(to); err == nil {
				t.Errorf("terminal state %s can move to %s", from, to)
			}
		}
	}
}

// Invalid values are a different failure from a valid move that is not
// allowed now. The error must not claim to be a transition error.
func TestTurnStateTransitionToRejectsInvalidValues(t *testing.T) {
	invalid := []TurnState{"", "running", "pending", "Accepted", " generating", "complete\n"}

	for _, bad := range invalid {
		for _, good := range allTurnStates {
			cases := map[string]func() error{
				fmt.Sprintf("from_%q_to_%s", string(bad), good): func() error {
					return bad.TransitionTo(good)
				},
				fmt.Sprintf("from_%s_to_%q", good, string(bad)): func() error {
					return good.TransitionTo(bad)
				},
			}

			for name, call := range cases {
				t.Run(name, func(t *testing.T) {
					err := call()
					if err == nil {
						t.Fatal("transition with an invalid state value allowed")
					}

					if errors.Is(err, ErrInvalidTurnTransition) {
						t.Fatalf("invalid value reported as a transition error: %v", err)
					}
				})
			}
		}
	}
}

// Graph properties checked through the implementation: every state is
// reachable from accepted, and every state can still reach a terminal state.
func TestTurnStateGraphHasNoUnreachableOrDeadEndStates(t *testing.T) {
	reachable := map[TurnState]bool{TurnStateAccepted: true}
	queue := []TurnState{TurnStateAccepted}

	for len(queue) > 0 {
		from := queue[0]
		queue = queue[1:]

		for _, to := range allTurnStates {
			if from.TransitionTo(to) == nil && !reachable[to] {
				reachable[to] = true
				queue = append(queue, to)
			}
		}
	}

	for _, state := range allTurnStates {
		if !reachable[state] {
			t.Errorf("%s is unreachable from accepted", state)
		}
	}

	for _, start := range allTurnStates {
		seen := map[TurnState]bool{start: true}
		stack := []TurnState{start}
		endsTerminal := false

		for len(stack) > 0 && !endsTerminal {
			from := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			if from.IsTerminal() {
				endsTerminal = true
				break
			}

			for _, to := range allTurnStates {
				if from.TransitionTo(to) == nil && !seen[to] {
					seen[to] = true
					stack = append(stack, to)
				}
			}
		}

		if !endsTerminal {
			t.Errorf("%s cannot reach any terminal state", start)
		}
	}
}

// The happy path for a full-history turn, step by step.
func TestFullHistoryTurnPath(t *testing.T) {
	path := []TurnState{
		TurnStateAccepted,
		TurnStateGenerating,
		TurnStateAnswerSaved,
		TurnStateComplete,
	}

	for i := 0; i+1 < len(path); i++ {
		if err := path[i].TransitionTo(path[i+1]); err != nil {
			t.Fatalf("step %s -> %s: %v", path[i], path[i+1], err)
		}
	}
}
