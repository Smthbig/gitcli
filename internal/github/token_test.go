package github

import "testing"

func TestGetTokenPrefersEnvironment(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv(EnvTokenName, "env-token")

	if err := Save("file-token"); err != nil {
		t.Fatalf("save token: %v", err)
	}

	if got := GetToken(); got != "env-token" {
		t.Fatalf("GetToken() = %q, want env-token", got)
	}
	if source := TokenSource(); source != "environment" {
		t.Fatalf("TokenSource() = %q, want environment", source)
	}
	if !HasStoredToken() {
		t.Fatalf("expected stored token to remain available")
	}
}
