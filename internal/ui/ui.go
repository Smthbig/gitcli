package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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
	fmt.Print(Cyan + label + ": " + Reset)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func SecretInput(label string) string {
	fmt.Print(Cyan + label + ": " + Reset)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
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
		fmt.Print(Yellow + question + " (" + defStr + "): " + Reset)
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
	fmt.Println(Cyan + label + Reset)
	for i, opt := range options {
		fmt.Printf(" %d) %s\n", i+1, opt)
	}

	for {
		fmt.Print("Select option: ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		for i := range options {
			if text == fmt.Sprint(i+1) {
				return i + 1
			}
		}
		fmt.Println(Red + "Invalid choice" + Reset)
	}
}

/* ============================================================
   Screen helpers
   ============================================================ */

func Pause() {
	fmt.Print("\nPress Enter to continue...")
	reader.ReadString('\n')
}

func Clear() {
	fmt.Print("\033[H\033[2J")
}

func Header(title string) {
	fmt.Println(Magenta + "========================================" + Reset)
	fmt.Println(Bold + Cyan + " " + title + Reset)
	fmt.Println(Magenta + "========================================" + Reset)
}

func Divider() {
	fmt.Println(Magenta + "----------------------------------------" + Reset)
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

func PrintKV(key, value string) {
	fmt.Printf("%-10s : %s\n", key, value)
}

func KeyHint(keys string) {
	fmt.Println(Blue + "[" + keys + "]" + Reset)
}

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
