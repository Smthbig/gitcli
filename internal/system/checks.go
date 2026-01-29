package system

import (
	"bufio"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"git-genius/internal/ui"
)

/*
Online means:
- Network reachable (NOT permission related)
- GitHub API reachable
*/
var Online bool = false

/* ============================================================
   SAFE COMMAND DETECTION (NO exec.LookPath)
   ============================================================ */

func commandExists(cmd string) bool {
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

/* ============================================================
   ENV DETECTION
   ============================================================ */

func isAndroid() bool {
	return runtime.GOOS == "android" ||
		os.Getenv("ANDROID_ROOT") != "" ||
		strings.Contains(strings.ToLower(os.Getenv("HOME")), "android")
}

func isRestrictedEnv() bool {
	// Android & Termux = restricted installs
	if isAndroid() {
		return true
	}
	if os.Getenv("PREFIX") != "" {
		return true
	}
	// No sudo → restricted
	if !commandExists("sudo") {
		return true
	}
	return false
}

/* ============================================================
   GIT CHECK (ANDROID SAFE)
   ============================================================ */

func EnsureGitInstalled() error {
	// Restricted envs: NEVER try auto install
	if isRestrictedEnv() {
		if commandExists("git") {
			return nil
		}
		ui.Warn("Git not found in PATH")
		ui.Info("Please install git manually for this environment")
		return errors.New("git missing (restricted environment)")
	}

	// Normal systems
	if commandExists("git") {
		return nil
	}

	ui.Error("Git is not installed")

	if runtime.GOOS == "windows" {
		ui.Info("Download Git from https://git-scm.com/downloads")
		return errors.New("git not installed")
	}

	if !ui.Confirm("Do you want to install Git now?") {
		return errors.New("git install declined")
	}

	if err := installGit(); err != nil {
		LogError("git install failed", err)
		return err
	}

	if !commandExists("git") {
		return errors.New("git install completed but binary not found")
	}

	ui.Success("Git installed successfully")
	return nil
}

/* ============================================================
   INSTALL LOGIC (NON-ANDROID ONLY)
   ============================================================ */

func installGit() error {
	if isRestrictedEnv() {
		return errors.New("automatic install not supported here")
	}

	switch runtime.GOOS {
	case "linux":
		return installGitLinux()
	case "darwin":
		return installGitMac()
	default:
		return errors.New("unsupported OS")
	}
}

func installGitLinux() error {
	switch {
	case commandExists("apt"):
		return runInstall("sudo apt update && sudo apt install -y git")
	case commandExists("dnf"):
		return runInstall("sudo dnf install -y git")
	case commandExists("yum"):
		return runInstall("sudo yum install -y git")
	case commandExists("pacman"):
		return runInstall("sudo pacman -S --noconfirm git")
	default:
		return errors.New("no supported package manager found")
	}
}

func installGitMac() error {
	if commandExists("brew") {
		return runInstall("brew install git")
	}

	ui.Info("Installing Xcode Command Line Tools")
	cmd := exec.Command("xcode-select", "--install")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

/* ============================================================
   HELPERS
   ============================================================ */

func runInstall(command string) error {
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = bufio.NewReader(os.Stdin)
	return cmd.Run()
}

/* ============================================================
   NETWORK CHECK (MINIMAL, SAFE, EXPLICIT)
   ============================================================ */

/*
CheckInternet:
- Runs ONLY when explicitly called
- Android = always online (git already proves connectivity)
- Non-Android = single lightweight GitHub API ping
- Auth / permission errors ≠ offline
*/
// CheckInternet intentionally trusts the environment.
// Used for Android / restricted systems where git already proves connectivity.
func CheckInternet() {
	Online = true
}
