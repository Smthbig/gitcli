package gitops

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"git-genius/internal/config"
)

func runGit(t *testing.T, dir string, args ...string) string {
	t.Helper()

	out, err := runGitAllowError(dir, args...)
	if err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, out)
	}
	return out
}

func runGitAllowError(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	if dir != "" {
		cmd.Dir = dir
	}

	out, err := cmd.CombinedOutput()
	return string(out), err
}

func withWorkingDir(t *testing.T, dir string) {
	t.Helper()

	prev, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir(%s): %v", dir, err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(prev)
	})
}

func TestPushAlsoPushesExistingLocalCommitsWhenWorkingTreeIsClean(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}

	home := t.TempDir()
	localParent, err := os.MkdirTemp("/tmp", "gitgenius-local-*")
	if err != nil {
		t.Fatalf("mktemp local: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(localParent) })
	remoteParent, err := os.MkdirTemp("/tmp", "gitgenius-remote-*")
	if err != nil {
		t.Fatalf("mktemp remote: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(remoteParent) })

	localRepo := filepath.Join(localParent, "repo")
	remoteRepo := filepath.Join(remoteParent, "remote.git")

	t.Setenv("HOME", home)

	runGit(t, "", "init", "--bare", remoteRepo)
	runGit(t, "", "init", localRepo)
	runGit(t, localRepo, "config", "user.name", "Test User")
	runGit(t, localRepo, "config", "user.email", "test@example.com")

	withWorkingDir(t, localRepo)

	if err := os.WriteFile(filepath.Join(localRepo, "README.md"), []byte("one\n"), 0o644); err != nil {
		t.Fatalf("write README: %v", err)
	}
	runGit(t, localRepo, "add", "README.md")
	runGit(t, localRepo, "commit", "-m", "initial")
	runGit(t, localRepo, "branch", "-M", "main")
	runGit(t, localRepo, "remote", "add", "origin", remoteRepo)
	initialPushOut, err := runGitAllowError(localRepo, "push", "-u", "origin", "main")
	if err != nil {
		if strings.Contains(initialPushOut, "bad pack") || strings.Contains(initialPushOut, "remote rejected") {
			t.Skip("sandbox git transport rejects local bare-remote pushes in this environment")
		}
		t.Fatalf("initial push failed: %v\n%s", err, initialPushOut)
	}

	if err := os.WriteFile(filepath.Join(localRepo, "README.md"), []byte("one\ntwo\n"), 0o644); err != nil {
		t.Fatalf("update README: %v", err)
	}
	runGit(t, localRepo, "add", "README.md")
	runGit(t, localRepo, "commit", "-m", "second")

	config.Save(config.Config{
		Branch:  "main",
		Remote:  "origin",
		WorkDir: localRepo,
	})

	if ok := Push(""); !ok {
		t.Fatalf("Push should push existing local commits on a clean working tree")
	}

	localHead := runGit(t, localRepo, "rev-parse", "HEAD")
	remoteHead := runGit(t, "", "--git-dir", remoteRepo, "rev-parse", "refs/heads/main")
	if localHead != remoteHead {
		t.Fatalf("remote HEAD %q does not match local HEAD %q", remoteHead, localHead)
	}
}
