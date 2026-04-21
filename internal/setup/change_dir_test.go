package setup

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"git-genius/internal/config"
)

func runGitSetupTest(t *testing.T, dir string, args ...string) {
	t.Helper()

	cmd := exec.Command("git", args...)
	if dir != "" {
		cmd.Dir = dir
	}
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, string(out))
	}
}

func initSetupRepo(t *testing.T, parent, name string) string {
	t.Helper()

	repo := filepath.Join(parent, name)
	runGitSetupTest(t, "", "init", repo)
	runGitSetupTest(t, repo, "config", "user.name", "Test User")
	runGitSetupTest(t, repo, "config", "user.email", "test@example.com")
	return repo
}

func TestActivateProjectDirLoadsTargetRepoConfigWithoutOverwrite(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}

	home := t.TempDir()
	parent := t.TempDir()
	t.Setenv("HOME", home)

	repoA := initSetupRepo(t, parent, "repo-a")
	repoB := initSetupRepo(t, parent, "repo-b")

	config.Save(config.Config{
		Branch:  "feature/a",
		Remote:  "origin",
		Owner:   "octo",
		Repo:    "repo-a",
		WorkDir: repoA,
	})
	config.Save(config.Config{
		Branch:  "release/b",
		Remote:  "upstream",
		Owner:   "octo",
		Repo:    "repo-b",
		WorkDir: repoB,
	})

	config.SetActiveWorkDir(repoA)

	got, isRepo, err := activateProjectDir(repoB)
	if err != nil {
		t.Fatalf("activateProjectDir: %v", err)
	}
	if !isRepo {
		t.Fatalf("expected repoB to be recognized as a git repo")
	}
	if got.WorkDir != repoB || got.Branch != "release/b" || got.Remote != "upstream" || got.Repo != "repo-b" {
		t.Fatalf("unexpected activated repo config: %+v", got)
	}

	if active := config.Load(); active.WorkDir != repoB || active.Branch != "release/b" || active.Remote != "upstream" {
		t.Fatalf("active config mismatch after switch: %+v", active)
	}

	if repoAConfig := config.LoadForWorkDir(repoA); repoAConfig.Branch != "feature/a" || repoAConfig.Repo != "repo-a" {
		t.Fatalf("repoA config should remain unchanged: %+v", repoAConfig)
	}
}

func TestActivateProjectDirReturnsDefaultConfigForPlainDirectory(t *testing.T) {
	home := t.TempDir()
	dir := t.TempDir()
	t.Setenv("HOME", home)

	got, isRepo, err := activateProjectDir(dir)
	if err != nil {
		t.Fatalf("activateProjectDir plain dir: %v", err)
	}
	if isRepo {
		t.Fatalf("did not expect plain directory to be treated as git repo")
	}
	if got.WorkDir != dir {
		t.Fatalf("workdir mismatch: got %q want %q", got.WorkDir, dir)
	}
	if config.Load().WorkDir != dir {
		t.Fatalf("expected active workdir to be updated to %q", dir)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Fatalf("plain directory should still exist: %v", err)
	}
}
