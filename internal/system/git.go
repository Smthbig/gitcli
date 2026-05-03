package system

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"

	"git-genius/internal/config"
	"git-genius/internal/github"
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

	// Always prepend --no-pager to avoid environment-specific pager failures
	fullArgs := append([]string{"--no-pager"}, args...)
	cmd := exec.Command(git, fullArgs...)

	cfg := config.Load()
	if cfg.WorkDir != "" {
		cmd.Dir = cfg.WorkDir
	}

	applyGitHubAuth(cmd, "")

	return cmd, nil
}

func GitCmdWithRemote(remote string, args ...string) (*exec.Cmd, error) {
	cmd, err := GitCmd(args...)
	if err != nil {
		return nil, err
	}

	applyGitHubAuth(cmd, remote)
	return cmd, nil
}

func GitCmdAt(dir string, args ...string) (*exec.Cmd, error) {
	git, err := getGitPath()
	if err != nil {
		return nil, err
	}

	// Always prepend --no-pager to avoid environment-specific pager failures
	fullArgs := append([]string{"--no-pager"}, args...)
	cmd := exec.Command(git, fullArgs...)
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

func RunGitWithRemote(remote string, args ...string) error {
	cmd, err := GitCmdWithRemote(remote, args...)
	if err != nil {
		return err
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func RunGitWithRemoteBuffered(remote string, args ...string) (string, error) {
	cmd, err := GitCmdWithRemote(remote, args...)
	if err != nil {
		return "", err
	}

	var stderr bytes.Buffer
	cmd.Stdout = os.Stdout
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderr)
	cmd.Stdin = os.Stdin

	err = cmd.Run()
	return strings.TrimSpace(stderr.String()), err
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

func GitOutputWithRemote(remote string, args ...string) (string, error) {
	cmd, err := GitCmdWithRemote(remote, args...)
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
// SHARED HELPERS
///////////////////////////////////////////////////////////////

func splitNonEmptyLines(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}

	lines := strings.Split(raw, "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			out = append(out, line)
		}
	}
	return out
}

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

func LocalBranches() ([]string, error) {
	out, err := GitOutput("for-each-ref", "--format=%(refname:short)", "refs/heads")
	if err != nil {
		return nil, err
	}

	branches := splitNonEmptyLines(out)
	sort.Strings(branches)
	return branches, nil
}

func HasLocalBranch(name string) bool {
	if strings.TrimSpace(name) == "" {
		return false
	}

	branches, err := LocalBranches()
	if err != nil {
		return false
	}

	for _, branch := range branches {
		if branch == name {
			return true
		}
	}
	return false
}

func SwitchToBranch(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("branch name is required")
	}
	if !HasLocalBranch(name) {
		return errors.New("branch does not exist: " + name)
	}
	return RunGit("checkout", name)
}

func CreateBranch(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("branch name is required")
	}
	if HasLocalBranch(name) {
		return errors.New("branch already exists: " + name)
	}
	return RunGit("checkout", "-b", name)
}

func PrepareBranch(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil
	}

	current := CurrentGitBranch()
	if current == name {
		return nil
	}

	if HasLocalBranch(name) {
		return SwitchToBranch(name)
	}

	return CreateBranch(name)
}

func RemoteNames() ([]string, error) {
	out, err := GitOutput("remote")
	if err != nil {
		return nil, err
	}

	remotes := splitNonEmptyLines(out)
	sort.Strings(remotes)
	return remotes, nil
}

func RemoteURL(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", errors.New("remote name is required")
	}
	return GitOutput("remote", "get-url", name)
}

func HasRemote(name string) bool {
	_, err := RemoteURL(name)
	return err == nil
}

func EnsureRemote(name, url string) error {
	name = strings.TrimSpace(name)
	url = strings.TrimSpace(url)
	if name == "" || url == "" {
		return errors.New("remote name and URL are required")
	}

	if HasRemote(name) {
		return RunGit("remote", "set-url", name, url)
	}

	return RunGit("remote", "add", name, url)
}

func HasRemoteTrackingBranch(remote, branch string) bool {
	remote = strings.TrimSpace(remote)
	branch = strings.TrimSpace(branch)
	if remote == "" || branch == "" {
		return false
	}

	_, err := GitOutput("rev-parse", "--verify", remote+"/"+branch)
	return err == nil
}

func AheadBehind(localRef, remoteRef string) (int, int, error) {
	localRef = strings.TrimSpace(localRef)
	remoteRef = strings.TrimSpace(remoteRef)
	if localRef == "" || remoteRef == "" {
		return 0, 0, errors.New("local and remote refs are required")
	}

	out, err := GitOutput("rev-list", "--left-right", "--count", localRef+"..."+remoteRef)
	if err != nil {
		return 0, 0, err
	}

	fields := strings.Fields(out)
	if len(fields) != 2 {
		return 0, 0, errors.New("unexpected ahead/behind output")
	}

	ahead, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, 0, err
	}
	behind, err := strconv.Atoi(fields[1])
	if err != nil {
		return 0, 0, err
	}

	return ahead, behind, nil
}

func RemoteUsesHTTPS(name string) bool {
	url, err := RemoteURL(name)
	if err != nil {
		return false
	}

	url = strings.TrimSpace(strings.ToLower(url))
	return strings.HasPrefix(url, "https://") || strings.HasPrefix(url, "http://")
}

func GitCredentialHelper() (string, error) {
	out, err := GitOutput("config", "--global", "--get", "credential.helper")
	if err != nil {
		return "", nil
	}
	return strings.TrimSpace(out), nil
}

func HasGitCredentialHelper() bool {
	helper, _ := GitCredentialHelper()
	return helper != ""
}

func SetGitCredentialHelper(value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return errors.New("credential helper value is required")
	}
	return RunGit("config", "--global", "credential.helper", value)
}

func ClearGitCredentialHelper() error {
	helper, _ := GitCredentialHelper()
	if helper == "" {
		return nil
	}
	return RunGit("config", "--global", "--unset-all", "credential.helper")
}

func ApproveGitHubCredential(username, token string) error {
	username = strings.TrimSpace(username)
	token = strings.TrimSpace(token)
	if username == "" || token == "" {
		return errors.New("username and token are required")
	}

	helper, _ := GitCredentialHelper()
	if helper == "" {
		return errors.New("git credential helper is not configured")
	}

	git, err := getGitPath()
	if err != nil {
		return err
	}

	cmd := exec.Command(git, "credential", "approve")
	cmd.Stdin = strings.NewReader(
		"protocol=https\n" +
			"host=github.com\n" +
			"username=" + username + "\n" +
			"password=" + token + "\n\n",
	)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if stderr.Len() > 0 {
			return errors.New(strings.TrimSpace(stderr.String()))
		}
		return err
	}

	return nil
}

func applyGitHubAuth(cmd *exec.Cmd, remote string) {
	if cmd == nil {
		return
	}

	token := github.GetToken()
	if token == "" {
		return
	}

	if !shouldUseGitHubAuth(remote) {
		return
	}

	username := resolveGitHubUsername()
	if username == "" {
		return
	}

	env := envMap(cmd.Env)
	env["GIT_TERMINAL_PROMPT"] = "0"
	env["GIT_GENIUS_GITHUB_USERNAME"] = username
	env["GIT_GENIUS_GITHUB_TOKEN"] = token
	env["GIT_CONFIG_COUNT"] = "1"
	env["GIT_CONFIG_KEY_0"] = "credential.helper"
	env["GIT_CONFIG_VALUE_0"] = `!f() { test "$1" = get || exit 0; printf '%s\n' "username=$GIT_GENIUS_GITHUB_USERNAME" "password=$GIT_GENIUS_GITHUB_TOKEN"; }; f`
	cmd.Env = envSlice(env)
}

func shouldUseGitHubAuth(remote string) bool {
	remote = strings.TrimSpace(remote)
	if remote == "" {
		return false
	}

	url, err := RemoteURL(remote)
	if err != nil {
		return false
	}

	return isGitHubHTTPSURL(url)
}

func isGitHubHTTPSURL(raw string) bool {
	raw = strings.TrimSpace(raw)
	return strings.HasPrefix(raw, "https://github.com/") || strings.HasPrefix(raw, "http://github.com/")
}

func resolveGitHubUsername() string {
	if username := github.GetUsername(); username != "" {
		return username
	}

	cfg := config.Load()
	if cfg.Owner != "" {
		return cfg.Owner
	}

	return "x-access-token"
}

func envMap(base []string) map[string]string {
	env := make(map[string]string, len(os.Environ())+len(base))
	for _, entry := range os.Environ() {
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) == 2 {
			env[parts[0]] = parts[1]
		}
	}
	for _, entry := range base {
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) == 2 {
			env[parts[0]] = parts[1]
		}
	}
	return env
}

func envSlice(env map[string]string) []string {
	keys := make([]string, 0, len(env))
	for key := range env {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	out := make([]string, 0, len(keys))
	for _, key := range keys {
		out = append(out, key+"="+env[key])
	}
	return out
}

func EnsureBranchSync() bool {
	if !IsGitRepo() {
		return true
	}

	cfg := config.Load()
	current := CurrentGitBranch()

	if current == "" {
		return true
	}

	if cfg.Branch == "" {
		cfg.Branch = current
		config.Save(cfg)
		return true
	}

	if current == cfg.Branch {
		return true
	}

	ui.Warn("Branch mismatch detected")
	ui.Info("Configured branch : " + cfg.Branch)
	ui.Info("Git branch        : " + current)

	if ui.ConfirmDefault("Use current git branch in config?", true) {
		cfg.Branch = current
		config.Save(cfg)
		ui.Success("Config branch synced to: " + current)
		return true
	}

	if ui.Confirm("Rename git branch to " + cfg.Branch + "?") {
		if err := RunGit("branch", "-m", cfg.Branch); err == nil {
			ui.Success("Git branch renamed to: " + cfg.Branch)
			return true
		}
		ui.Error("Failed to rename git branch")
		return false
	}

	ui.Warn("Keeping current branch and syncing config to stay consistent")
	cfg.Branch = current
	config.Save(cfg)
	return true
}
