package system

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"git-genius/internal/config"
	"git-genius/internal/ui"
)

var (
	gitPath string
	gitErr  error
	gitOnce sync.Once
)

///////////////////////////////////////////////////////////////
// COMMAND EXISTS (RESTORED FOR DOCTOR)
///////////////////////////////////////////////////////////////

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

///////////////////////////////////////////////////////////////
// GIT RESOLUTION
///////////////////////////////////////////////////////////////

func resolveGit() (string, error) {
	pathEnv := os.Getenv("PATH")
	if pathEnv == "" {
		return "", errors.New("PATH not set")
	}

	for _, dir := range strings.Split(pathEnv, ":") {
		full := filepath.Join(dir, "git")
		info, err := os.Stat(full)
		if err == nil && info.Mode().IsRegular() && info.Mode()&0111 != 0 {
			return full, nil
		}
	}

	return "", errors.New("git not found in PATH")
}

func getGitPath() (string, error) {
	gitOnce.Do(func() {
		gitPath, gitErr = resolveGit()
	})

	if gitPath == "" {
		return "", errors.New("git not installed")
	}

	return gitPath, gitErr
}

///////////////////////////////////////////////////////////////
// SAFE DIRECTORY (Git ≥2.35 Android fix)
///////////////////////////////////////////////////////////////

func EnsureSafeDirectory(path string) {
	if path == "" {
		return
	}

	// prevent duplicate entries
	existing, _ := GitOutput("config", "--global", "--get-all", "safe.directory")
	if strings.Contains(existing, path) {
		return
	}

	_ = RunGit("config", "--global", "--add", "safe.directory", path)
}

///////////////////////////////////////////////////////////////
// COMMAND BUILDERS
///////////////////////////////////////////////////////////////

func GitCmd(args ...string) (*exec.Cmd, error) {
	git, err := getGitPath()
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(git, args...)

	cfg := config.Load()
	if cfg.WorkDir != "" {
		cmd.Dir = cfg.WorkDir
	}

	return cmd, nil
}

func GitCmdAt(dir string, args ...string) (*exec.Cmd, error) {
	git, err := getGitPath()
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(git, args...)
	cmd.Dir = dir
	return cmd, nil
}

///////////////////////////////////////////////////////////////
// EXECUTORS
///////////////////////////////////////////////////////////////

func RunGit(args ...string) error {
	cmd, err := GitCmd(args...)
	if err != nil {
		return err
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func RunGitAt(dir string, args ...string) error {
	cmd, err := GitCmdAt(dir, args...)
	if err != nil {
		return err
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func GitOutput(args ...string) (string, error) {
	cmd, err := GitCmd(args...)
	if err != nil {
		return "", err
	}

	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		if stderr.Len() > 0 {
			return "", errors.New(strings.TrimSpace(stderr.String()))
		}
		return "", err
	}

	return strings.TrimSpace(out.String()), nil
}

func GitOutputAt(dir string, args ...string) (string, error) {
	cmd, err := GitCmdAt(dir, args...)
	if err != nil {
		return "", err
	}

	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		if stderr.Len() > 0 {
			return "", errors.New(strings.TrimSpace(stderr.String()))
		}
		return "", err
	}

	return strings.TrimSpace(out.String()), nil
}

///////////////////////////////////////////////////////////////
// REPOSITORY CHECKS
///////////////////////////////////////////////////////////////

func IsGitRepo() bool {
	cmd, err := GitCmd("rev-parse", "--is-inside-work-tree")
	if err != nil {
		return false
	}
	return cmd.Run() == nil
}

func IsGitRepoAt(dir string) bool {
	cmd, err := GitCmdAt(dir, "rev-parse", "--is-inside-work-tree")
	if err != nil {
		return false
	}
	return cmd.Run() == nil
}

func EnsureGitRepo() bool {
	cfg := config.Load()
	EnsureSafeDirectory(cfg.WorkDir)

	if IsGitRepo() {
		return true
	}

	ui.Warn("Selected directory is not a git repository")

	if !ui.Confirm("Initialize a git repository here?") {
		ui.Error("Git repository required")
		return false
	}

	if err := RunGit("init"); err != nil {
		ui.Error("Git initialization failed")
		return false
	}

	ui.Success("Git repository initialized")
	return true
}

///////////////////////////////////////////////////////////////
// BRANCH HELPERS
///////////////////////////////////////////////////////////////

func CurrentGitBranch() string {
	branch, err := GitOutput("branch", "--show-current")
	if err != nil {
		return ""
	}
	return branch
}

func CurrentRemote() string {
	cfg := config.Load()

	if cfg.Remote == "" {
		return "-"
	}

	// verify remote actually exists
	_, err := GitOutput("remote", "get-url", cfg.Remote)
	if err != nil {
		return "-"
	}

	return cfg.Remote
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

	if current == "" || current == cfg.Branch {
		return true
	}

	ui.Warn("Branch mismatch detected")
	ui.Info("Configured branch : " + cfg.Branch)
	ui.Info("Git branch        : " + current)

	if ui.Confirm("Rename git branch to " + cfg.Branch + "?") {
		if err := RunGit("branch", "-m", cfg.Branch); err == nil {
			ui.Success("Git branch renamed")
			return true
		}
		ui.Warn("Branch rename failed")
	}

	cfg.Branch = current
	config.Save(cfg)
	return true
}
