package main

import (
	"slices"
	"testing"
)

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

func TestModelAllowlistIDsReturnsSortedConfiguredModels(t *testing.T) {
	want := []string{"local-model-a", "local-model-b", "local-model-c"}
	for _, input := range [][]string{
		{"local-model-c", "local-model-a", "local-model-b"},
		{"local-model-b", "local-model-c", "local-model-a"},
	} {
		allowlist, err := NewModelAllowlist(input)
		if err != nil {
			t.Fatal(err)
		}
		got := allowlist.IDs()
		if !slices.Equal(got, want) {
			t.Fatalf("IDs = %v, want %v", got, want)
		}
		for _, id := range got {
			if !allowlist.Allows(id) {
				t.Errorf("listed model %q is not allowed", id)
			}
		}
		if allowlist.Allows("local-model-unconfigured") {
			t.Error("listing changed allowed membership")
		}
	}
}

func TestModelAllowlistIDsReturnsIndependentSlice(t *testing.T) {
	allowlist, err := NewModelAllowlist([]string{"local-model-b", "local-model-a"})
	if err != nil {
		t.Fatal(err)
	}
	first, second := allowlist.IDs(), allowlist.IDs()
	if len(first) != 2 || len(second) != 2 {
		t.Fatal("expected two configured models")
	}
	first[0] = "synthetic-injected-model"
	want := []string{"local-model-a", "local-model-b"}
	if !slices.Equal(second, want) || !slices.Equal(allowlist.IDs(), want) {
		t.Error("mutating one result changed another or subsequent result")
	}
	if !allowlist.Allows("local-model-a") || allowlist.Allows("synthetic-injected-model") {
		t.Error("mutating IDs changed authorization")
	}
}

func TestZeroModelAllowlistIDsReturnsEmptySlice(t *testing.T) {
	var allowlist ModelAllowlist
	if got := allowlist.IDs(); got == nil || len(got) != 0 {
		t.Fatalf("zero-value IDs = %#v, want initialized empty slice", got)
	}
}
