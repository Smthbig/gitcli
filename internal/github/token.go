package github

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

const (
	geniusDirName = ".git-genius"
	tokenFileName = "token"
)

/*
getTokenPath:
  - Stores token in user home directory
  - Example:
    ~/.git-genius/token
*/
func getTokenPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(home, geniusDirName)
	return filepath.Join(dir, tokenFileName), nil
}

/* ================= TOKEN ================= */

func GetToken() string {
	path, err := getTokenPath()
	if err != nil {
		return ""
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(data))
}

func Save(token string) error {
	token = strings.TrimSpace(token)
	if token == "" {
		return errors.New("empty token")
	}

	path, err := getTokenPath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)

	// Create ~/.git-genius
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	// Write token securely
	return os.WriteFile(path, []byte(token), 0600)
}

func Delete() error {
	path, err := getTokenPath()
	if err != nil {
		return err
	}
	return os.Remove(path)
}
