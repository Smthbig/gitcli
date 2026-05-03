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
   ANSI color codes (kept for legacy support if needed)
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
	fmt.Print(InfoStyle.Render(label) + " ❯ ")
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func InputDefault(label, defaultValue string) string {
	if strings.TrimSpace(defaultValue) == "" {
		return Input(label)
	}

	fmt.Print(InfoStyle.Render(label) + DimStyle.Render(" ["+defaultValue+"]") + " ❯ ")
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
		fmt.Print(WarningStyle.Render(question) + DimStyle.Render(" ("+defStr+")") + " ❯ ")
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
		fmt.Println(ErrorStyle.Render("Please enter y or n."))
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
		fmt.Println(ErrorStyle.Render("Invalid choice"))
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
	fmt.Print("\n" + WarningStyle.Render("Press Enter to continue..."))
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
	fmt.Println(HeaderStyle.Render("⚡ GIT GENIUS ⚡"))
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
	leftWidth := 46

	if termWidth < 95 {
		for _, l := range leftCol { fmt.Println(l) }
		fmt.Println()
		for _, r := range rightCol { fmt.Println(r) }
		return
	}

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

		fmt.Printf("%s   %s\n", left, right)
	}
}

func BoxHeader(title string) {
	fmt.Println(HeaderStyle.Render(title))
}

func BoxMenu(title string, options []string) {
	var builder strings.Builder
	builder.WriteString(TitleStyle.Render(" " + title + " ") + "\n\n")
	for _, opt := range options {
		builder.WriteString("  " + opt + "\n")
	}
	fmt.Println(MainBorderStyle.Render(builder.String()))
}

func Info(msg string) {
	fmt.Println(InfoStyle.Render("ℹ " + msg))
}

func Success(msg string) {
	fmt.Println(SuccessStyle.Render("✔ " + msg))
}

func Warn(msg string) {
	fmt.Println(WarningStyle.Render("⚠ " + msg))
}

func Error(msg string) {
	fmt.Println("\a" + ErrorStyle.Render("✘ " + msg))
}

func Divider() {
	fmt.Println(DividerStyle.Render(strings.Repeat("─", 44)))
}

/* ============================================================
   UX Helpers
   ============================================================ */

func Spinner(message string) func() {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	stop := make(chan bool)
	
	fmt.Printf(InfoStyle.Render(" " + frames[0] + " ") + message + "... ")
	
	go func() {
		i := 0
		for {
			select {
			case <-stop:
				return
			default:
				fmt.Printf("\r"+InfoStyle.Render(" %s ")+message+"... ", frames[i])
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

func OpenURL(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "android":
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
   Help Renderer
   ============================================================ */

func PrintHelp(lines []string) {
	for _, line := range lines {
		fmt.Println("  " + line)
	}
	fmt.Println()
}
