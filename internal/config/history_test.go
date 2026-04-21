package config

import (
	"strings"
	"testing"
)

func containsMessage(messages []string, needle string) bool {
	for _, message := range messages {
		if strings.Contains(message, needle) {
			return true
		}
	}
	return false
}

func TestHistorySuggestionsReflectRecentBehavior(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	workDir := t.TempDir()

	RecordHistory(workDir, "daily", "push", false, "push failed")
	RecordHistory(workDir, "daily", "push", true, "")
	RecordHistory(workDir, "daily", "push", true, "")
	RecordHistory(workDir, "daily", "push", true, "")
	RecordHistory(workDir, "tools", "change_project_dir", true, "")
	RecordHistory(workDir, "tools", "change_project_dir", true, "")

	got := HistorySuggestions(workDir)

	if !containsMessage(got, "Doctor first") {
		t.Fatalf("expected doctor suggestion, got %v", got)
	}
	if !containsMessage(got, "Smart Pull before push") {
		t.Fatalf("expected smart pull suggestion, got %v", got)
	}
	if !containsMessage(got, "Change Project Directory") {
		t.Fatalf("expected recent directory suggestion, got %v", got)
	}
}

func TestHasHistoryForWorkDirOnlyMatchesThatRepo(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	workDirA := t.TempDir()
	workDirB := t.TempDir()

	RecordHistory(workDirA, "daily", "status", true, "")

	if !HasHistoryForWorkDir(workDirA) {
		t.Fatalf("expected history for workDirA")
	}
	if HasHistoryForWorkDir(workDirB) {
		t.Fatalf("did not expect history for workDirB")
	}
}
