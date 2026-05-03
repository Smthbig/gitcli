package ui

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"golang.org/x/term"
)

/* ============================================================
   ANSI color codes
   ============================================================ */

const (
	Reset   = "\033[0m"
	Bold    = "\033[1m"
	Red     = "\033[1;31m"
	Green   = "\033[1;32m"
	Yellow  = "\033[1;33m"
	Blue    = "\033[1;34m"
	Cyan    = "\033[1;36m"
	Magenta = "\033[1;35m"
)

/* ============================================================
   Reader (single instance)
   ============================================================ */

var reader = bufio.NewReader(os.Stdin)

/* ============================================================
   Input helpers
   ============================================================ */

func Input(label string) string {
	fmt.Print(Cyan + label + Reset + " ❯ ")
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func InputDefault(label, defaultValue string) string {
	if strings.TrimSpace(defaultValue) == "" {
		return Input(label)
	}

	fmt.Print(Cyan + label + Reset + " [" + defaultValue + "] ❯ ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	if text == "" {
		return defaultValue
	}
	return text
}

func Confirm(question string) bool {
	return ConfirmDefault(question, false)
}

func ConfirmDefault(question string, def bool) bool {
	defStr := "y/N"
	if def {
		defStr = "Y/n"
	}

	for {
		fmt.Print(Yellow + question + Reset + " (" + defStr + ") ❯ ")
		text, _ := reader.ReadString('\n')
		ans := strings.ToLower(strings.TrimSpace(text))

		if ans == "" {
			return def
		}
		if ans == "y" || ans == "yes" {
			return true
		}
		if ans == "n" || ans == "no" {
			return false
		}
		fmt.Println(Red + "Please enter y or n." + Reset)
	}
}

/*
Select presents numbered options and returns choice index (1-based)
*/
func Select(label string, options []string) int {
	BoxMenu(label, options)

	for {
		text := Input("Select option")
		for i := range options {
			if text == fmt.Sprint(i+1) {
				return i + 1
			}
		}
		fmt.Println(Red + "Invalid choice" + Reset)
	}
}

func MenuChoice() string {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		// Fallback to regular input if raw mode fails
		return strings.ToLower(Input("Select option"))
	}
	defer func() {
		_ = term.Restore(int(os.Stdin.Fd()), oldState)
	}()

	b := make([]byte, 1)
	_, err = os.Stdin.Read(b)
	if err != nil {
		return ""
	}

	char := string(b)
	// Handle Ctrl+C (ASCII 3)
	if b[0] == 3 {
		fmt.Println("^C")
		os.Exit(0)
	}

	fmt.Println(char) // Echo the choice
	return strings.ToLower(char)
}

/* ============================================================
   Screen helpers
   ============================================================ */

func Pause() {
	fmt.Print("\n" + Yellow + "Press Enter to continue..." + Reset)
	reader.ReadString('\n')
}

func Clear() {
	fmt.Print("\033[H\033[2J")
}

func BoxHeader(title string) {
	width := 42
	fmt.Println(Magenta + "┏" + strings.Repeat("━", width-2) + "┓" + Reset)
	
	// Calculate padding for centering
	titleLen := len(title)
	padding := (width - 2 - titleLen) / 2
	if padding < 0 {
		padding = 0
	}
	rightPadding := width - 2 - titleLen - padding
	if rightPadding < 0 {
		rightPadding = 0
	}

	fmt.Printf(Magenta+"┃"+Reset+strings.Repeat(" ", padding)+Bold+Cyan+"%s"+Reset+strings.Repeat(" ", rightPadding)+Magenta+"┃\n"+Reset, title)
	fmt.Println(Magenta + "┗" + strings.Repeat("━", width-2) + "┛" + Reset)
}

func BoxMenu(title string, options []string) {
	width := 42
	fmt.Println(Magenta + "┏" + strings.Repeat("━", width-2) + "┓" + Reset)
	fmt.Printf(Magenta+"┃ "+Reset+Bold+Cyan+"%-*s"+Magenta+" ┃\n"+Reset, width-4, title)
	fmt.Println(Magenta + "┣" + strings.Repeat("━", width-2) + "┫" + Reset)
	
	for _, opt := range options {
		fmt.Printf(Magenta+"┃ "+Reset+"%-*s"+Magenta+" ┃\n"+Reset, width-4, opt)
	}
	fmt.Println(Magenta + "┗" + strings.Repeat("━", width-2) + "┛" + Reset)
}

func Divider() {
	fmt.Println(Magenta + "----------------------------------------" + Reset)
}

/* ============================================================
   External Helpers
   ============================================================ */

// OpenURL opens the specified URL in the default browser/app.
// Optimized for Termux via termux-open, but supports other OSs.
func OpenURL(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "android":
		// Termux specific
		cmd = exec.Command("termux-open", url)
	case "linux":
		if _, err := exec.LookPath("termux-open"); err == nil {
			cmd = exec.Command("termux-open", url)
		} else {
			cmd = exec.Command("xdg-open", url)
		}
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return cmd.Start()
}

func SelectConventionalType() string {
	types := []string{
		"feat:     A new feature",
		"fix:      A bug fix",
		"docs:     Documentation only changes",
		"style:    Changes that do not affect the meaning of the code",
		"refactor: A code change that neither fixes a bug nor adds a feature",
		"perf:     A code change that improves performance",
		"test:     Adding missing tests or correcting existing tests",
		"chore:    Changes to the build process or auxiliary tools",
	}

	choice := Select("Select change type", types)
	selected := types[choice-1]
	return strings.Split(selected, ":")[0]
}

/* ============================================================
   Message helpers
   ============================================================ */

func Info(msg string) {
	fmt.Println(Cyan + "ℹ " + msg + Reset)
}

func Success(msg string) {
	fmt.Println(Green + "✔ " + msg + Reset)
}

func Warn(msg string) {
	fmt.Println(Yellow + "⚠ " + msg + Reset)
}

func Error(msg string) {
	fmt.Println(Red + "✘ " + msg + Reset)
}

/* ============================================================
   Utility helpers
   ============================================================ */

// Help renders a help screen with title and bullet points
// ============================================================
// Help Renderer
// ============================================================

// PrintHelp prints help lines in a clean readable format
func PrintHelp(lines []string) {
	for _, line := range lines {
		fmt.Println("  " + line)
	}
	fmt.Println()
}
