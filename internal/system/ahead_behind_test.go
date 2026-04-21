package system

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestAheadBehindCountsLocalAndRemoteDifference(t *testing.T) {
	repo := initTestRepo(t)

	bare := filepath.Join(t.TempDir(), "remote.git")
	runGitCommand(t, "", "init", "--bare", bare)
	runGitCommand(t, repo, "remote", "add", "origin", bare)
	runGitCommand(t, repo, "branch", "-M", "main")
	out, err := exec.Command("git", "-C", repo, "push", "-u", "origin", "main").CombinedOutput()
	if err != nil {
		text := string(out)
		if strings.Contains(text, "bad pack") || strings.Contains(text, "remote rejected") {
			t.Skip("sandbox git transport rejects local bare-remote pushes in this environment")
		}
		t.Fatalf("git push failed: %v\n%s", err, text)
	}

	ahead, behind, err := AheadBehind("main", "origin/main")
	if err != nil || ahead != 0 || behind != 0 {
		t.Fatalf("AheadBehind initial = %d/%d err=%v", ahead, behind, err)
	}

	if err := os.WriteFile(filepath.Join(repo, "local.txt"), []byte("local change\n"), 0o644); err != nil {
		t.Fatalf("write local change: %v", err)
	}
	runGitCommand(t, repo, "add", "local.txt")
	runGitCommand(t, repo, "commit", "-m", "local change")

	ahead, behind, err = AheadBehind("main", "origin/main")
	if err != nil || ahead != 1 || behind != 0 {
		t.Fatalf("AheadBehind after local commit = %d/%d err=%v", ahead, behind, err)
	}
}
