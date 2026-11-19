package main

import "testing"

// The eight canonical turn-state names from docs/contracts.md#turn-states.
// They are written as string literals on purpose: the test must compile and
// report each missing state by name before the constants exist.
var canonicalTurnStates = []string{
	"accepted",
	"generating",
	"answer_saved",
	"memory_pending",
	"complete",
	"failed",
	"cancelled",
	"recovery_required",
}

func TestTurnStateValidateAcceptsCanonicalValues(t *testing.T) {
	if len(canonicalTurnStates) != 8 {
		t.Fatalf("fixture must list 8 canonical states, got %d",
			len(canonicalTurnStates))
	}

	for _, name := range canonicalTurnStates {
		t.Run(name, func(t *testing.T) {
			if err := TurnState(name).Validate(); err != nil {
				t.Fatalf("TurnState(%q).Validate() = %v, want nil", name, err)
			}
		})
	}
}

func TestTurnStateValidateRejectsNonCanonicalValues(t *testing.T) {
	tests := []struct {
		name  string
		state string
	}{
		{"empty", ""},
		{"space_only", " "},

		// Aliases the contract explicitly forbids or that look plausible.
		{"alias_running", "running"},
		{"alias_pending", "pending"},
		{"alias_completed", "completed"},
		{"alias_done", "done"},
		{"alias_canceled_us_spelling", "canceled"},
		{"alias_error", "error"},
		{"alias_saved", "saved"},

		// Session states use a separate type.
		{"session_state_deleting", "deleting"},
		{"session_state_deleted", "deleted"},

		// Case changes.
		{"upper_case", "ACCEPTED"},
		{"title_case", "Generating"},
		{"mixed_case", "answer_Saved"},

		// Padding and embedded whitespace.
		{"leading_space", " generating"},
		{"trailing_space", "generating "},
		{"both_spaces", " generating "},
		{"trailing_newline", "complete\n"},
		{"trailing_crlf", "complete\r\n"},
		{"leading_tab", "\tfailed"},
		{"trailing_nul", "failed\x00"},

		// Separator variants of multi-word names.
		{"hyphen_separator", "answer-saved"},
		{"space_separator", "memory pending"},
		{"no_separator", "recoveryrequired"},
		{"camel_case", "answerSaved"},

		// Near misses.
		{"prefix_only", "accept"},
		{"extra_suffix", "accepteds"},
		{"two_states_joined", "accepted,generating"},
		{"full_width_letters", "ａｃｃｅｐｔｅｄ"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := TurnState(test.state).Validate(); err == nil {
				t.Fatalf("TurnState(%q).Validate() = nil, want error",
					test.state)
			}
		})
	}
}

func TestTurnStateValidateRejectsPaddedAndRecasedCanonicalValues(t *testing.T) {
	for _, name := range canonicalTurnStates {
		variants := map[string]string{
			"leading_space":  " " + name,
			"trailing_space": name + " ",
			"upper_first":    string(name[0]-'a'+'A') + name[1:],
		}

		for label, variant := range variants {
			t.Run(name+"/"+label, func(t *testing.T) {
				if err := TurnState(variant).Validate(); err == nil {
					t.Fatalf("TurnState(%q).Validate() = nil, want error",
						variant)
				}
			})
		}
	}
}
