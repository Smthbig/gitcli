package config

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type HistoryEntry struct {
	Timestamp int64  `json:"timestamp"`
	WorkDir   string `json:"work_dir"`
	Section   string `json:"section"`
	Action    string `json:"action"`
	Success   bool   `json:"success"`
	Note      string `json:"note,omitempty"`
}

func RecordHistory(workDir, section, action string, success bool, note string) {
	p := historyPath()
	if p == "" {
		return
	}
	_ = os.MkdirAll(filepath.Dir(p), 0700)

	e := HistoryEntry{
		Timestamp: time.Now().UnixMilli(),
		WorkDir:   workDir,
		Section:   section,
		Action:    action,
		Success:   success,
		Note:      strings.TrimSpace(note),
	}

	b, err := json.Marshal(e)
	if err != nil {
		return
	}

	f, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return
	}
	defer f.Close()
	_, _ = f.Write(append(b, '\n'))
}

func RecentHistory(limit int) []HistoryEntry {
	if limit <= 0 {
		limit = 20
	}
	p := historyPath()
	if p == "" {
		return nil
	}
	f, err := os.Open(p)
	if err != nil {
		return nil
	}
	defer f.Close()

	var all []HistoryEntry
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" {
			continue
		}
		var e HistoryEntry
		if json.Unmarshal([]byte(line), &e) == nil {
			all = append(all, e)
		}
	}
	if len(all) <= limit {
		return all
	}
	return all[len(all)-limit:]
}

func HistorySuggestions(workDir string) []string {
	recent := RecentHistory(60)
	if len(recent) == 0 {
		return []string{"Run Tools -> Setup / Reconfigure to initialize this project."}
	}

	var inDir []HistoryEntry
	for _, e := range recent {
		if e.WorkDir == workDir {
			inDir = append(inDir, e)
		}
	}
	if len(inDir) == 0 {
		return []string{"No history for this repo yet. Start with Setup / Reconfigure."}
	}

	failures := 0
	pushCount := 0
	pullCount := 0
	changeDirCount := 0
	for _, e := range inDir {
		if !e.Success {
			failures++
		}
		if e.Action == "push" {
			pushCount++
		}
		if e.Action == "pull" || e.Action == "smart_pull" {
			pullCount++
		}
		if e.Action == "change_project_dir" {
			changeDirCount++
		}
	}

	out := make([]string, 0, 3)
	if failures > 0 {
		out = append(out, "Recent failures detected. Run Tools -> Doctor first.")
	}
	if pushCount >= 3 && pullCount == 0 {
		out = append(out, "You push often without pulling. Use Smart Pull before push to avoid conflicts.")
	}
	if changeDirCount >= 2 {
		out = append(out, "You switch repos often. Use recent-directory shortcuts in Setup or Change Project Directory.")
	}
	if len(out) == 0 {
		out = append(out, "Workflow looks healthy. Continue with Daily Git Operations.")
	}
	return out
}

func historyPath() string {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return ""
	}
	return filepath.Join(home, ".git-genius", "history.jsonl")
}
