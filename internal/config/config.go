package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	Dir  = ".git/.genius"
	File = Dir + "/config.json"
)

const (
	stateDirName  = ".git-genius"
	stateFileName = "state.json"
)

type appState struct {
	ActiveWorkDir  string   `json:"active_work_dir"`
	RecentWorkDirs []string `json:"recent_work_dirs"`
}

// Config holds Git Genius configuration
type Config struct {
	/* ---------------- Git basics ---------------- */
	Branch string `json:"branch"`
	Remote string `json:"remote"`

	/* ---------------- GitHub repo ---------------- */
	Owner string `json:"owner"` // username or organisation
	Repo  string `json:"repo"`  // repository name

	/* ---------------- Push state ---------------- */
	FirstPushDone bool `json:"first_push_done"`

	/* ---------------- Project directory ---------------- */
	// Empty = current working directory
	WorkDir string `json:"work_dir"`
}

/* ============================================================
   Load / Save
   ============================================================ */

// Load reads config from disk and applies safe defaults
func Load() Config {
	return loadForWorkDir(preferredWorkDir(), true)
}

func LoadForWorkDir(workDir string) Config {
	return loadForWorkDir(workDir, false)
}

// Save writes config with secure permissions
func Save(c Config) {
	applyDefaults(&c)
	normalizePaths(&c)

	if c.WorkDir == "" {
		c.WorkDir = preferredWorkDir()
	}

	p := configPath(c.WorkDir)
	gitDir := filepath.Join(c.WorkDir, ".git")
	if info, err := os.Stat(gitDir); err == nil && info.IsDir() {
		_ = os.MkdirAll(filepath.Dir(p), 0700)

		data, err := json.MarshalIndent(c, "", "  ")
		if err == nil {
			_ = os.WriteFile(p, data, 0600)
		}
	}

	updateState(c.WorkDir)
}

func SetActiveWorkDir(workDir string) {
	if workDir == "" {
		return
	}
	updateState(workDir)
}

/* ============================================================
   Defaults & helpers
   ============================================================ */

func defaultConfig() Config {
	return Config{
		Branch:        "main",
		Remote:        "origin",
		Owner:         "",
		Repo:          "",
		FirstPushDone: false,
		WorkDir:       "",
	}
}

// applyDefaults keeps backward compatibility
func applyDefaults(c *Config) {
	if c.Branch == "" {
		c.Branch = "main"
	}
	if c.Remote == "" {
		c.Remote = "origin"
	}
}

// normalizePaths ensures WorkDir is absolute
func normalizePaths(c *Config) {
	if c.WorkDir == "" {
		return
	}

	abs, err := filepath.Abs(c.WorkDir)
	if err == nil {
		c.WorkDir = abs
	}
}

// GetWorkDir returns resolved working directory
func (c Config) GetWorkDir() string {
	if c.WorkDir != "" {
		return c.WorkDir
	}
	if d := preferredWorkDir(); d != "" {
		return d
	}
	wd, _ := os.Getwd()
	return wd
}

func RecentWorkDirs() []string {
	s := loadState()
	out := make([]string, 0, len(s.RecentWorkDirs))
	for _, d := range s.RecentWorkDirs {
		if d != "" {
			if info, err := os.Stat(d); err != nil || !info.IsDir() {
				continue
			}
			out = append(out, d)
		}
	}
	return out
}

func preferredWorkDir() string {
	s := loadState()
	if s.ActiveWorkDir != "" {
		if info, err := os.Stat(s.ActiveWorkDir); err == nil && info.IsDir() {
			return s.ActiveWorkDir
		}
	}
	wd, _ := os.Getwd()
	return wd
}

func configPath(workDir string) string {
	if workDir == "" {
		return File
	}
	return filepath.Join(workDir, ".git", ".genius", "config.json")
}

func loadForWorkDir(workDir string, remember bool) Config {
	if workDir == "" {
		workDir = preferredWorkDir()
	}

	data, err := os.ReadFile(configPath(workDir))
	if err != nil {
		c := defaultConfig()
		c.WorkDir = workDir
		if remember {
			updateState(c.WorkDir)
		}
		return c
	}

	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		c = defaultConfig()
		c.WorkDir = workDir
		if remember {
			updateState(c.WorkDir)
		}
		return c
	}

	applyDefaults(&c)
	normalizePaths(&c)
	if c.WorkDir == "" {
		c.WorkDir = workDir
	}
	if remember {
		updateState(c.WorkDir)
	}

	return c
}

func HasProjectConfig(workDir string) bool {
	if workDir == "" {
		workDir = preferredWorkDir()
	}

	info, err := os.Stat(configPath(workDir))
	return err == nil && !info.IsDir()
}

func statePath() string {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return ""
	}
	return filepath.Join(home, stateDirName, stateFileName)
}

func loadState() appState {
	p := statePath()
	if p == "" {
		return appState{}
	}

	data, err := os.ReadFile(p)
	if err != nil {
		return appState{}
	}

	var s appState
	if err := json.Unmarshal(data, &s); err != nil {
		return appState{}
	}
	return s
}

func saveState(s appState) {
	p := statePath()
	if p == "" {
		return
	}

	_ = os.MkdirAll(filepath.Dir(p), 0700)
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(p, b, 0600)
}

func updateState(workDir string) {
	if workDir == "" {
		return
	}

	abs, err := filepath.Abs(workDir)
	if err != nil {
		return
	}

	s := loadState()
	s.ActiveWorkDir = abs

	next := []string{abs}
	for _, d := range s.RecentWorkDirs {
		if d != abs && d != "" {
			next = append(next, d)
		}
		if len(next) >= 10 {
			break
		}
	}
	s.RecentWorkDirs = next
	saveState(s)
}
