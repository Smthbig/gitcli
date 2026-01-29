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

// Config holds Git Genius configuration
type Config struct {
	/* ---------------- Git basics ---------------- */
	Branch        string `json:"branch"`
	DefaultBranch string `json:"default_branch"` // main / master
	Remote        string `json:"remote"`

	/* ---------------- GitHub repo ---------------- */
	Owner       string `json:"owner"`        // username or organisation
	Repo        string `json:"repo"`         // repository name
	IsOrgRepo   bool   `json:"is_org_repo"`  // user vs org
	OrgName     string `json:"org_name"`     // organisation name
	PrivateRepo bool   `json:"private_repo"` // public / private
	RepoCreated bool   `json:"repo_created"` // GitHub repo exists or not

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
	data, err := os.ReadFile(File)
	if err != nil {
		return defaultConfig()
	}

	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return defaultConfig()
	}

	applyDefaults(&c)
	normalizePaths(&c)

	return c
}

// Save writes config with secure permissions
func Save(c Config) {
	applyDefaults(&c)
	normalizePaths(&c)

	_ = os.MkdirAll(Dir, 0700)

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return
	}

	_ = os.WriteFile(File, data, 0600)
}

/* ============================================================
   Defaults & helpers
   ============================================================ */

func defaultConfig() Config {
	return Config{
		Branch:        "main",
		DefaultBranch: "main",
		Remote:        "origin",
		Owner:         "",
		Repo:          "",
		IsOrgRepo:     false,
		OrgName:       "",
		PrivateRepo:   false,
		RepoCreated:   false,
		FirstPushDone: false,
		WorkDir:       "",
	}
}

// applyDefaults keeps backward compatibility
func applyDefaults(c *Config) {
	if c.Branch == "" {
		c.Branch = "main"
	}
	if c.DefaultBranch == "" {
		c.DefaultBranch = c.Branch
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
	wd, _ := os.Getwd()
	return wd
}
