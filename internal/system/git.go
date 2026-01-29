package system

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"git-genius/internal/config"
	"git-genius/internal/ui"
)

//
// ============================================================
// SAFE COMMAND / GIT RESOLUTION (NO exec.LookPath)
// ============================================================
//

var (
	gitPath string
	gitOnce sync.Once
)

// CommandExists checks any command safely (Android safe)
func CommandExists(cmd string) bool {
	pathEnv := os.Getenv("PATH")
	if pathEnv == "" {
		return false
	}

	for _, dir := range strings.Split(pathEnv, ":") {
		full := filepath.Join(dir, cmd)
		info, err := os.Stat(full)
		if err == nil && info.Mode().IsRegular() && info.Mode()&0111 != 0 {
			return true
		}
	}
	return false
}

func resolveGit() string {
	pathEnv := os.Getenv("PATH")
	if pathEnv == "" {
		return ""
	}

	for _, dir := range strings.Split(pathEnv, ":") {
		full := filepath.Join(dir, "git")
		info, err := os.Stat(full)
		if err == nil && info.Mode().IsRegular() && info.Mode()&0111 != 0 {
			return full
		}
	}
	return ""
}

func getGitPath() string {
	gitOnce.Do(func() {
		gitPath = resolveGit()
	})
	return gitPath
}

//
// ============================================================
// SAFE DIRECTORY (Git ≥ 2.35 ANDROID FIX)
// ============================================================
//

// EnsureSafeDirectory fixes "dubious ownership" issue
func EnsureSafeDirectory(dir string) {
	if dir == "" {
		return
	}

	_ = exec.Command(
		getGitPath(),
		"config",
		"--global",
		"--add",
		"safe.directory",
		dir,
	).Run()
}

//
// ============================================================
// GIT COMMAND BUILDERS
// ============================================================
//

// GitCmd builds git command using config.WorkDir
func GitCmd(args ...string) *exec.Cmd {
	git := getGitPath()
	if git == "" {
		return exec.Command("false")
	}

	cmd := exec.Command(git, args...)

	cfg := config.Load()
	if cfg.WorkDir != "" {
		cmd.Dir = cfg.WorkDir
	}

	return cmd
}

// GitCmdAt builds git command for explicit directory
func GitCmdAt(dir string, args ...string) *exec.Cmd {
	git := getGitPath()
	if git == "" {
		return exec.Command("false")
	}

	cmd := exec.Command(git, args...)
	cmd.Dir = dir
	return cmd
}

//
// ============================================================
// GIT EXECUTORS
// ============================================================
//

// RunGit runs git in config.WorkDir
func RunGit(args ...string) error {
	cmd := GitCmd(args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		LogError("git "+strings.Join(args, " "), err)
		return err
	}
	return nil
}

// RunGitAt runs git in specific directory
func RunGitAt(dir string, args ...string) error {
	cmd := GitCmdAt(dir, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		LogError("git "+strings.Join(args, " "), err)
		return err
	}
	return nil
}

// GitOutput runs git and returns trimmed output
func GitOutput(args ...string) (string, error) {
	cmd := GitCmd(args...)
	out, err := cmd.Output()
	if err != nil {
		LogError("git "+strings.Join(args, " "), err)
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// GitOutputAt runs git in dir and returns output
func GitOutputAt(dir string, args ...string) (string, error) {
	cmd := GitCmdAt(dir, args...)
	out, err := cmd.Output()
	if err != nil {
		LogError("git "+strings.Join(args, " "), err)
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

//
// ============================================================
// REPOSITORY CHECKS
// ============================================================
//

func IsGitRepo() bool {
	cmd := GitCmd("rev-parse", "--is-inside-work-tree")
	return cmd.Run() == nil
}

func IsGitRepoAt(dir string) bool {
	cmd := GitCmdAt(dir, "rev-parse", "--is-inside-work-tree")
	return cmd.Run() == nil
}

func EnsureGitRepo() bool {
	cfg := config.Load()
	EnsureSafeDirectory(cfg.WorkDir)

	if IsGitRepo() {
		return true
	}

	ui.Warn("Selected directory is not a git repository")

	if !ui.Confirm("Do you want to initialize a git repository here?") {
		ui.Error("Git repository required to continue")
		return false
	}

	if err := RunGit("init"); err != nil {
		ui.Error("Failed to initialize git repository")
		return false
	}

	ui.Success("Git repository initialized")
	return true
}

func EnsureGitRepoAt(dir string) bool {
	EnsureSafeDirectory(dir)

	if IsGitRepoAt(dir) {
		return true
	}

	ui.Warn("Selected directory is not a git repository")

	if !ui.Confirm("Do you want to initialize a git repository here?") {
		ui.Error("Git repository required to continue")
		return false
	}

	if err := RunGitAt(dir, "init"); err != nil {
		ui.Error("Failed to initialize git repository")
		return false
	}

	ui.Success("Git repository initialized")
	return true
}

//
// ============================================================
// BRANCH HELPERS
// ============================================================
//

func CurrentGitBranch() string {
	branch, err := GitOutput("branch", "--show-current")
	if err != nil {
		return ""
	}
	return branch
}

func CurrentGitBranchAt(dir string) string {
	branch, err := GitOutputAt(dir, "branch", "--show-current")
	if err != nil {
		return ""
	}
	return branch
}

func EnsureBranchSync() bool {
	cfg := config.Load()
	current := CurrentGitBranch()

	if current == "" || cfg.Branch == current {
		return true
	}

	ui.Warn("Branch mismatch detected")
	ui.Info("Configured branch : " + cfg.Branch)
	ui.Info("Git branch        : " + current)

	if ui.Confirm("Rename git branch to " + cfg.Branch + "?") {
		if err := RunGit("branch", "-m", cfg.Branch); err == nil {
			ui.Success("Git branch renamed to: " + cfg.Branch)
			return true
		}
		ui.Warn("Branch rename failed")
	}

	cfg.Branch = current
	config.Save(cfg)
	ui.Success("Config branch updated to: " + current)
	return true
}
