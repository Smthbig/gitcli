package ui

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

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

func GetTermWidth() int {
	width, _, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		return 80 // Default fallback
	}
	return width
}

func Logo() {
	fmt.Println(Bold + Cyan + `
   ⚡ GIT GENIUS ⚡
   PRO-TERMINAL v2` + Reset)
}

func CyberSparkline(data []int) string {
	if len(data) == 0 {
		return ""
	}
	// Unicode blocks for a sparkline effect
	bars := []string{" ", "▂", "▃", "▄", "▅", "▆", "▇", "█"}
	max := 0
	for _, v := range data {
		if v > max {
			max = v
		}
	}
	if max == 0 {
		max = 1
	}

	var res strings.Builder
	for _, v := range data {
		idx := (v * (len(bars) - 1)) / max
		res.WriteString(bars[idx])
	}
	return res.String()
}

func RenderGrid(leftCol, rightCol []string) {
	termWidth := GetTermWidth()
	// We'll use 46 as the fixed width for our components (+ borders/padding)
	leftWidth := 46

	if termWidth < 95 {
		// Stack them if screen is narrow (Termux portrait)
		for _, l := range leftCol { fmt.Println(l) }
		fmt.Println()
		for _, r := range rightCol { fmt.Println(r) }
		return
	}

	// Side-by-side (Termux landscape or Tablet)
	maxLines := len(leftCol)
	if len(rightCol) > maxLines {
		maxLines = len(rightCol)
	}

	for i := 0; i < maxLines; i++ {
		left := ""
		if i < len(leftCol) {
			left = leftCol[i]
		} else {
			left = strings.Repeat(" ", leftWidth)
		}
		
		right := ""
		if i < len(rightCol) {
			right = rightCol[i]
		}

		// Use a safe way to print columns with ANSI escapes
		// Note: This is an experimental layout approach
		fmt.Printf("%s   %s\n", left, right)
	}
}

func BoxHeader(title string) {
	width := 44
	fmt.Println(Magenta + "╔" + strings.Repeat("═", width-2) + "╗" + Reset)
	
	titleLen := len(title)
	padding := (width - 2 - titleLen) / 2
	if padding < 0 { padding = 0 }
	rightPadding := width - 2 - titleLen - padding
	if rightPadding < 0 { rightPadding = 0 }

	fmt.Printf(Magenta+"║"+Reset+strings.Repeat(" ", padding)+Bold+Cyan+"%s"+Reset+strings.Repeat(" ", rightPadding)+Magenta+"║\n"+Reset, title)
	fmt.Println(Magenta + "╚" + strings.Repeat("═", width-2) + "╝" + Reset)
}

func BoxMenu(title string, options []string) {
	for _, l := range GetBoxMenuLines(title, options) {
		fmt.Println(l)
	}
}

func GetBoxMenuLines(title string, options []string) []string {
	width := 44
	var res []string
	res = append(res, Cyan+"╔"+strings.Repeat("═", width-2)+"╗"+Reset)
	res = append(res, fmt.Sprintf(Cyan+"║ "+Reset+Bold+Yellow+"%-*s"+Cyan+" ║"+Reset, width-4, title))
	res = append(res, Cyan+"╠"+strings.Repeat("═", width-2)+"╣"+Reset)
	
	for _, opt := range options {
		if opt == "" {
			res = append(res, Cyan+"║"+strings.Repeat(" ", width-2)+"║"+Reset)
			continue
		}
		res = append(res, fmt.Sprintf(Cyan+"║ "+Reset+"%-*s"+Cyan+" ║"+Reset, width-4, opt))
	}
	res = append(res, Cyan+"╚"+strings.Repeat("═", width-2)+"╝"+Reset)
	return res
}

func BoxPanel(title string, lines []string, color string) {
	for _, l := range GetBoxPanelLines(title, lines, color) {
		fmt.Println(l)
	}
}

func GetBoxPanelLines(title string, lines []string, color string) []string {
	width := 44
	if color == "" { color = Blue }
	var res []string
	res = append(res, color+"┌"+strings.Repeat("─", width-2)+"┐"+Reset)
	if title != "" {
		res = append(res, fmt.Sprintf(color+"│ "+Reset+Bold+title+strings.Repeat(" ", width-len(title)-3)+color+"│"+Reset))
		res = append(res, color+"├"+strings.Repeat("─", width-2)+"┤"+Reset)
	}
	
	for _, line := range lines {
		displayLine := line
		if len(line) > width-4 {
			displayLine = line[:width-7] + "..."
		}
		res = append(res, fmt.Sprintf(color+"│ "+Reset+"%-*s"+color+" │"+Reset, width-4, displayLine))
	}
	res = append(res, color+"└"+strings.Repeat("─", width-2)+"┘"+Reset)
	return res
}

func Divider() {
	fmt.Println(Magenta + "----------------------------------------" + Reset)
}

/* ============================================================
   UX Helpers
   ============================================================ */

// Spinner shows an animation while a task is running.
// Returns a stop function.
func Spinner(message string) func() {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	stop := make(chan bool)
	
	fmt.Printf(Cyan+" %s "+Reset+"%s... ", frames[0], message)
	
	go func() {
		i := 0
		for {
			select {
			case <-stop:
				return
			default:
				fmt.Printf("\r"+Cyan+" %s "+Reset+"%s... ", frames[i], message)
				i = (i + 1) % len(frames)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	return func() {
		stop <- true
		fmt.Print("\r" + strings.Repeat(" ", len(message)+10) + "\r") // Clear line
	}
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
	fmt.Println("\a" + Red + "✘ " + msg + Reset)
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
