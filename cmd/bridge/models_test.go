package main

import "testing"

func TestModelAllowlistAllowsConfiguredModelsOnly(t *testing.T) {
	allowlist, err := NewModelAllowlist([]string{
		"local-model-a",
		"local-model-b",
	})
	if err != nil {
		t.Fatalf("NewModelAllowlist returned an unexpected error: %v", err)
	}

	tests := []struct {
		name    string
		modelID string
		want    bool
	}{
		{name: "first configured model", modelID: "local-model-a", want: true},
		{name: "second configured model", modelID: "local-model-b", want: true},
		{name: "unknown model", modelID: "local-model-c", want: false},
		{name: "case changed", modelID: "LOCAL-MODEL-A", want: false},
		{name: "surrounding spaces", modelID: " local-model-a ", want: false},
		{name: "empty model", modelID: "", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := allowlist.Allows(tt.modelID); got != tt.want {
				t.Errorf("Allows(%q) = %v, want %v", tt.modelID, got, tt.want)
			}
		})
	}
}

func TestNewModelAllowlistRejectsInvalidConfiguration(t *testing.T) {
	tests := []struct {
		name     string
		modelIDs []string
	}{
		{name: "no configured models"},
		{name: "blank model ID", modelIDs: []string{"local-model-a", "   "}},
		{name: "model ID with surrounding spaces", modelIDs: []string{" local-model-a "}},
		{name: "duplicate model ID", modelIDs: []string{"local-model-a", "local-model-a"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := NewModelAllowlist(tt.modelIDs); err == nil {
				t.Fatal("NewModelAllowlist returned nil error for invalid configuration")
			}
		})
	}
}

func TestModelAllowlistOwnsItsConfiguration(t *testing.T) {
	modelIDs := []string{"local-model-a"}

	allowlist, err := NewModelAllowlist(modelIDs)
	if err != nil {
		t.Fatalf("NewModelAllowlist returned an unexpected error: %v", err)
	}

	modelIDs[0] = "changed-after-construction"

	if !allowlist.Allows("local-model-a") {
		t.Fatal("allowlist changed when the caller mutated its input slice")
	}
	if allowlist.Allows("changed-after-construction") {
		t.Fatal("allowlist retained the caller's mutable slice")
	}
}
