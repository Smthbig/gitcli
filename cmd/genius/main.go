package main

import (
	"fmt"
	"os"

	"git-genius/internal/menu"
	"git-genius/internal/system"
	"git-genius/internal/ui"
)

var version = "dev"

func init() {
	//  ABSOLUTE SAFETY for Android / restricted kernels
	// Must be set before ANY syscall-heavy logic runs
	_ = os.Setenv("GODEBUG", "faccessat2=0")
}

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--version", "-v", "version":
			fmt.Println(version)
			return
		}
	}

	ui.Clear()
	ui.BoxHeader("Git Genius")

	// --- Git availability (safe & non-fatal on Android) ---
	gitAvailable := true
	if err := system.EnsureGitInstalled(); err != nil {
		gitAvailable = false
		ui.Warn("Git is not available")

		if system.IsRestrictedRuntime() {
			ui.Info("Running in restricted mode (Android / container)")
			ui.Info("Only setup, doctor, and help flows are available until Git is installed")
		} else {
			ui.Error("Git is required for full functionality")
			ui.Pause()
			return
		}
	}

	// --- Start UI ---
	menu.Start(version, gitAvailable)
}
