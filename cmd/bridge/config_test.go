package main

import (
	"net/http/httptest"
	"os"
	"slices"
	"strings"
	"testing"
)

func TestModelsFromEnvironmentLoadsConfiguredIDs(t *testing.T) {
	t.Setenv("BRIDGE_MODEL_IDS", "synthetic-model-b,synthetic-model-a")
	allowlist, err := modelsFromEnvironment()
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("BRIDGE_MODEL_IDS", "synthetic-model-replacement")
	want := []string{"synthetic-model-a", "synthetic-model-b"}
	if !slices.Equal(allowlist.IDs(), want) {
		t.Fatalf("IDs = %v, want %v", allowlist.IDs(), want)
	}
	for _, id := range want {
		if !allowlist.Allows(id) {
			t.Errorf("configured model %q denied", id)
		}
	}
	if allowlist.Allows("synthetic-model-replacement") {
		t.Error("later environment change altered allowlist")
	}
}

func TestModelsFromEnvironmentRejectsInvalidLists(t *testing.T) {
	for _, raw := range []string{"", " ", ",", "synthetic-model-a,", ",synthetic-model-a", "synthetic-model-a,,synthetic-model-b", "synthetic-model-a, synthetic-model-b", "synthetic-model-a,synthetic-model-a"} {
		t.Run(raw, func(t *testing.T) {
			t.Setenv("BRIDGE_MODEL_IDS", raw)
			allowlist, err := modelsFromEnvironment()
			if err == nil {
				t.Fatal("invalid model configuration accepted")
			}
			if len(allowlist.IDs()) != 0 || allowlist.Allows("synthetic-model-a") {
				t.Error("invalid configuration returned a partial allowlist")
			}
		})
	}
}

func TestModelsFromEnvironmentRejectsUnsetVariable(t *testing.T) {
	t.Setenv("BRIDGE_MODEL_IDS", "synthetic-temporary")
	if err := os.Unsetenv("BRIDGE_MODEL_IDS"); err != nil {
		t.Fatal(err)
	}
	if _, err := modelsFromEnvironment(); err == nil {
		t.Fatal("unset model configuration accepted")
	}
}

func TestAuthFromEnvironmentLoadsIndependentConfiguration(t *testing.T) {
	// Synthetic values override ambient configuration and are restored by testing.
	t.Setenv("BRIDGE_TOKEN", "synthetic-env-token")
	t.Setenv("BRIDGE_LAB_ID", "synthetic-lab")
	t.Setenv("BRIDGE_USER_ID", "synthetic-user")
	auth, err := authFromEnvironment()
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("BRIDGE_TOKEN", "synthetic-replacement-token")
	t.Setenv("BRIDGE_LAB_ID", "synthetic-other-lab")
	t.Setenv("BRIDGE_USER_ID", "synthetic-other-user")
	r := httptest.NewRequest("GET", "/v1/models", nil)
	r.Header.Set("Authorization", "Bearer synthetic-env-token")
	got, err := auth.Authenticate(r)
	want := Principal{LabID: "synthetic-lab", UserID: "synthetic-user"}
	if err != nil || got != want {
		t.Fatalf("identity = %+v, %v; want %+v", got, err, want)
	}
	r.Header.Set("Authorization", "Bearer synthetic-replacement-token")
	if _, err := auth.Authenticate(r); err == nil {
		t.Fatal("environment change silently changed existing authenticator")
	}
}

func TestAuthFromEnvironmentRejectsInvalidConfiguration(t *testing.T) {
	for _, tt := range []struct{ name, key, value string }{
		{"empty token", "BRIDGE_TOKEN", ""},
		{"blank token", "BRIDGE_TOKEN", " \t"},
		{"embedded whitespace", "BRIDGE_TOKEN", "synthetic-env token"},
		{"empty lab", "BRIDGE_LAB_ID", ""},
		{"blank lab", "BRIDGE_LAB_ID", " \t"},
		{"empty user", "BRIDGE_USER_ID", ""},
		{"blank user", "BRIDGE_USER_ID", " \t"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("BRIDGE_TOKEN", "synthetic-env-token")
			t.Setenv("BRIDGE_LAB_ID", "synthetic-lab")
			t.Setenv("BRIDGE_USER_ID", "synthetic-user")
			t.Setenv(tt.key, tt.value)
			auth, err := authFromEnvironment()
			if err == nil {
				t.Fatal("invalid configuration accepted")
			}
			if strings.Contains(err.Error(), "synthetic-") {
				t.Fatal("configuration error leaked configured values")
			}
			r := httptest.NewRequest("GET", "/v1/models", nil)
			r.Header.Set("Authorization", "Bearer synthetic-env-token")
			if _, authErr := auth.Authenticate(r); authErr == nil {
				t.Fatal("failed constructor returned a usable authenticator")
			}
		})
	}
}
