package config

import (
	"os"
	"path/filepath"
	"testing"
)

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

func TestSaveLoadAndHasProjectConfig(t *testing.T) {
	home := t.TempDir()
	repo := t.TempDir()
	if err := os.MkdirAll(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	t.Setenv("HOME", home)
	withWorkingDir(t, repo)

	want := Config{
		Branch:        "feature/setup",
		Remote:        "upstream",
		Owner:         "octo-org",
		Repo:          "git-genius",
		FirstPushDone: true,
		WorkDir:       repo,
	}

	Save(want)

	if !HasProjectConfig(repo) {
		t.Fatalf("expected repo config to be written for %s", repo)
	}

	got := Load()
	if got.Branch != want.Branch {
		t.Fatalf("branch mismatch: got %q want %q", got.Branch, want.Branch)
	}
	if got.Remote != want.Remote {
		t.Fatalf("remote mismatch: got %q want %q", got.Remote, want.Remote)
	}
	if got.Owner != want.Owner || got.Repo != want.Repo {
		t.Fatalf("repo mismatch: got %s/%s want %s/%s", got.Owner, got.Repo, want.Owner, want.Repo)
	}
	if got.WorkDir != repo {
		t.Fatalf("workdir mismatch: got %q want %q", got.WorkDir, repo)
	}

	recent := RecentWorkDirs()
	if len(recent) == 0 || recent[0] != repo {
		t.Fatalf("recent workdirs missing repo: %v", recent)
	}
}

func TestSaveSkipsRepoConfigOutsideGitRepo(t *testing.T) {
	home := t.TempDir()
	dir := t.TempDir()

	t.Setenv("HOME", home)
	withWorkingDir(t, dir)

	Save(Config{
		Branch:  "dev",
		Remote:  "origin",
		WorkDir: dir,
	})

	if HasProjectConfig(dir) {
		t.Fatalf("did not expect repo config outside git repo")
	}

	recent := RecentWorkDirs()
	if len(recent) == 0 || recent[0] != dir {
		t.Fatalf("expected state to keep active workdir, got %v", recent)
	}
}

func TestRecentWorkDirsReordersWhenSwitchingProjects(t *testing.T) {
	home := t.TempDir()
	repoA := t.TempDir()
	repoB := t.TempDir()

	if err := os.MkdirAll(filepath.Join(repoA, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir repoA .git: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(repoB, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir repoB .git: %v", err)
	}

	t.Setenv("HOME", home)

	Save(Config{WorkDir: repoA})
	Save(Config{WorkDir: repoB})
	Save(Config{WorkDir: repoA})

	recent := RecentWorkDirs()
	if len(recent) < 2 {
		t.Fatalf("expected at least two recent directories, got %v", recent)
	}
	if recent[0] != repoA || recent[1] != repoB {
		t.Fatalf("unexpected recent workdir order: %v", recent)
	}
	if got := preferredWorkDir(); got != repoA {
		t.Fatalf("preferredWorkDir = %q want %q", got, repoA)
	}
}

func TestLoadForWorkDirAndSetActiveWorkDirKeepRepoConfigsSeparated(t *testing.T) {
	home := t.TempDir()
	repoA := t.TempDir()
	repoB := t.TempDir()

	if err := os.MkdirAll(filepath.Join(repoA, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir repoA .git: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(repoB, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir repoB .git: %v", err)
	}

	t.Setenv("HOME", home)

	Save(Config{Branch: "feature/a", Remote: "origin", Owner: "octo", Repo: "repo-a", WorkDir: repoA})
	Save(Config{Branch: "release/b", Remote: "upstream", Owner: "octo", Repo: "repo-b", WorkDir: repoB})

	SetActiveWorkDir(repoA)
	if got := Load(); got.Branch != "feature/a" || got.Repo != "repo-a" {
		t.Fatalf("active repoA config mismatch: %+v", got)
	}

	if got := LoadForWorkDir(repoB); got.Branch != "release/b" || got.Remote != "upstream" || got.Repo != "repo-b" {
		t.Fatalf("LoadForWorkDir repoB mismatch: %+v", got)
	}

	SetActiveWorkDir(repoB)
	if got := Load(); got.Branch != "release/b" || got.Remote != "upstream" || got.Repo != "repo-b" {
		t.Fatalf("active repoB config mismatch: %+v", got)
	}

	if got := LoadForWorkDir(repoA); got.Branch != "feature/a" || got.Repo != "repo-a" {
		t.Fatalf("repoA config was altered unexpectedly: %+v", got)
	}
}
