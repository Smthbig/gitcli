package main

import (
	"os"

	"git-genius/internal/menu"
	"git-genius/internal/system"
	"git-genius/internal/ui"
)

func init() {
	//  ABSOLUTE SAFETY for Android / restricted kernels
	// Must be set before ANY syscall-heavy logic runs
	_ = os.Setenv("GODEBUG", "faccessat2=0")
}

func main() {
	ui.Clear()
	ui.Header("Git Genius")

	// --- Git availability (safe & non-fatal on Android) ---
	if err := system.EnsureGitInstalled(); err != nil {
		ui.Warn("Git is not available")
		ui.Info("Some features will be limited")

		if system.IsRestrictedRuntime() {
			ui.Info("Running in restricted mode (Android / container)")
		} else {
			ui.Error("Git is required for full functionality")
			ui.Pause()
			return
		}
	}

	// --- Network check (best-effort only) ---
	//	system.CheckInternet()

	// --- Start UI ---
	menu.Start()
}
