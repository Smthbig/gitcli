package gitops

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"git-genius/internal/config"
)

func createCommitInRepo(repo, name, contents, message string) error {
	if err := os.WriteFile(filepath.Join(repo, name), []byte(contents), 0o644); err != nil {
		return err
	}
	if _, err := runGitAllowError(repo, "add", name); err != nil {
		return err
	}
	_, err := runGitAllowError(repo, "commit", "-m", message)
	return err
}

func TestInspectRepoStateDetectsFirstPushAndAheadBehind(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}

	home := t.TempDir()
	localRepo := filepath.Join(t.TempDir(), "repo")
	remoteRepo := filepath.Join(t.TempDir(), "remote.git")

	t.Setenv("HOME", home)

	runGit(t, "", "init", "--bare", remoteRepo)
	runGit(t, "", "init", localRepo)
	runGit(t, localRepo, "config", "user.name", "Test User")
	runGit(t, localRepo, "config", "user.email", "test@example.com")

	withWorkingDir(t, localRepo)

	if err := createCommitInRepo(localRepo, "README.md", "one\n", "initial"); err != nil {
		t.Fatalf("create initial commit: %v", err)
	}
	runGit(t, localRepo, "branch", "-M", "main")
	runGit(t, localRepo, "remote", "add", "origin", remoteRepo)

	config.Save(config.Config{
		Branch:  "main",
		Remote:  "origin",
		WorkDir: localRepo,
	})

	state := InspectRepoState()
	if !state.NeedsFirstPush {
		t.Fatalf("expected first push to be pending: %+v", state)
	}

	pushOut, err := runGitAllowError(localRepo, "push", "-u", "origin", "main")
	if err != nil {
		if strings.Contains(pushOut, "bad pack") || strings.Contains(pushOut, "remote rejected") {
			t.Skip("sandbox git transport rejects local bare-remote pushes in this environment")
		}
		t.Fatalf("push failed: %v\n%s", err, pushOut)
	}

	state = InspectRepoState()
	if state.NeedsFirstPush {
		t.Fatalf("did not expect first push after publish: %+v", state)
	}
	if !state.HasAheadBehind || state.Ahead != 0 || state.Behind != 0 {
		t.Fatalf("expected synced ahead/behind after publish: %+v", state)
	}

	if err := createCommitInRepo(localRepo, "next.txt", "two\n", "second"); err != nil {
		t.Fatalf("create second commit: %v", err)
	}

	state = InspectRepoState()
	if !state.HasAheadBehind || state.Ahead != 1 || state.Behind != 0 {
		t.Fatalf("expected local ahead after second commit: %+v", state)
	}
}
