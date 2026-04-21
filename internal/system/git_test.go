package system

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"git-genius/internal/config"
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

func TestApproveGitHubCredentialWritesThroughConfiguredHelper(t *testing.T) {
	if !CommandExists("git") {
		t.Skip("git not installed")
	}

	home := t.TempDir()
	t.Setenv("HOME", home)

	if err := SetGitCredentialHelper("store"); err != nil {
		t.Fatalf("SetGitCredentialHelper: %v", err)
	}

	if err := ApproveGitHubCredential("algospider", "token-123"); err != nil {
		t.Fatalf("ApproveGitHubCredential: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(home, ".git-credentials"))
	if err != nil {
		t.Fatalf("read credentials: %v", err)
	}

	got := string(data)
	if got == "" || got != "https://algospider:token-123@github.com\n" {
		t.Fatalf("unexpected credential store contents: %q", got)
	}
}

func TestGitCmdWithRemoteInjectsGitHubHTTPSCredentials(t *testing.T) {
	repo := initTestRepo(t)
	runGitCommand(t, repo, "remote", "add", "origin", "https://github.com/example/project.git")

	config.Save(config.Config{
		Branch:  "main",
		Remote:  "origin",
		Owner:   "saved-owner",
		Repo:    "project",
		WorkDir: repo,
	})

	t.Setenv("GIT_GENIUS_GITHUB_TOKEN", "env-token")
	t.Setenv("GIT_GENIUS_GITHUB_USERNAME", "saved-user")

	cmd, err := GitCmdWithRemote("origin", "push", "-u", "origin", "main")
	if err != nil {
		t.Fatalf("GitCmdWithRemote: %v", err)
	}

	env := strings.Join(cmd.Env, "\n")
	for _, want := range []string{
		"GIT_TERMINAL_PROMPT=0",
		"GIT_GENIUS_GITHUB_TOKEN=env-token",
		"GIT_GENIUS_GITHUB_USERNAME=saved-user",
		"GIT_CONFIG_KEY_0=credential.helper",
	} {
		if !strings.Contains(env, want) {
			t.Fatalf("expected env to contain %q\n%s", want, env)
		}
	}
}

func TestGitCmdWithRemoteSkipsCredentialInjectionForNonGitHubRemote(t *testing.T) {
	repo := initTestRepo(t)
	runGitCommand(t, repo, "remote", "add", "origin", "https://example.com/project.git")

	config.Save(config.Config{
		Branch:  "main",
		Remote:  "origin",
		Owner:   "saved-owner",
		Repo:    "project",
		WorkDir: repo,
	})

	t.Setenv("GIT_GENIUS_GITHUB_TOKEN", "env-token")
	t.Setenv("GIT_GENIUS_GITHUB_USERNAME", "saved-user")

	cmd, err := GitCmdWithRemote("origin", "push", "-u", "origin", "main")
	if err != nil {
		t.Fatalf("GitCmdWithRemote: %v", err)
	}

	env := strings.Join(cmd.Env, "\n")
	if strings.Contains(env, "GIT_GENIUS_GITHUB_TOKEN=env-token") {
		t.Fatalf("did not expect GitHub credentials for non-GitHub remote\n%s", env)
	}
}
