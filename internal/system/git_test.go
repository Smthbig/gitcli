package system

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func runGitCommand(t *testing.T, dir string, args ...string) {
	t.Helper()

	cmd := exec.Command("git", args...)
	if dir != "" {
		cmd.Dir = dir
	}
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, string(out))
	}
}

func withRepoWorkingDir(t *testing.T, repo string) {
	t.Helper()

	prev, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(repo); err != nil {
		t.Fatalf("chdir(%s): %v", repo, err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(prev)
	})
}

func initTestRepo(t *testing.T) string {
	t.Helper()

	if !CommandExists("git") {
		t.Skip("git not installed")
	}

	t.Setenv("HOME", t.TempDir())
	repo := t.TempDir()
	runGitCommand(t, "", "init", repo)
	runGitCommand(t, repo, "config", "user.name", "Test User")
	runGitCommand(t, repo, "config", "user.email", "test@example.com")

	if err := os.WriteFile(filepath.Join(repo, "README.md"), []byte("hello\n"), 0o644); err != nil {
		t.Fatalf("write README: %v", err)
	}
	runGitCommand(t, repo, "add", "README.md")
	runGitCommand(t, repo, "commit", "-m", "initial")

	withRepoWorkingDir(t, repo)
	return repo
}

func TestEnsureRemoteUpdatesExistingAndAddsNewRemote(t *testing.T) {
	repo := initTestRepo(t)

	runGitCommand(t, repo, "remote", "add", "origin", "https://example.com/old.git")

	if err := EnsureRemote("origin", "https://example.com/new.git"); err != nil {
		t.Fatalf("EnsureRemote update failed: %v", err)
	}
	if got, err := GitOutput("remote", "get-url", "origin"); err != nil || got != "https://example.com/new.git" {
		t.Fatalf("origin URL = %q, err=%v", got, err)
	}

	if err := EnsureRemote("upstream", "https://example.com/upstream.git"); err != nil {
		t.Fatalf("EnsureRemote add failed: %v", err)
	}
	if got, err := GitOutput("remote", "get-url", "upstream"); err != nil || got != "https://example.com/upstream.git" {
		t.Fatalf("upstream URL = %q, err=%v", got, err)
	}
}

func TestCreateBranchAndPrepareBranchAreSafe(t *testing.T) {
	initTestRepo(t)

	initialBranch := CurrentGitBranch()
	if initialBranch == "" {
		t.Fatalf("expected initial branch after first commit")
	}

	if err := CreateBranch("feature/test"); err != nil {
		t.Fatalf("CreateBranch failed: %v", err)
	}
	if got := CurrentGitBranch(); got != "feature/test" {
		t.Fatalf("current branch = %q, want feature/test", got)
	}

	if err := SwitchToBranch(initialBranch); err != nil {
		t.Fatalf("SwitchToBranch(%q) failed: %v", initialBranch, err)
	}

	if err := CreateBranch("feature/test"); err == nil {
		t.Fatalf("expected CreateBranch to reject existing branch")
	}

	if err := PrepareBranch("release/test"); err != nil {
		t.Fatalf("PrepareBranch failed: %v", err)
	}
	if got := CurrentGitBranch(); got != "release/test" {
		t.Fatalf("current branch = %q, want release/test", got)
	}
}
