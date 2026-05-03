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
	userFileName  = "username"
	EnvTokenName  = "GIT_GENIUS_GITHUB_TOKEN"
	EnvUserName   = "GIT_GENIUS_GITHUB_USERNAME"
)

/*
getTokenPath:
  - Stores token in user home directory
  - Example:
    ~/.git-genius/token
*/
func getTokenPath() (string, error) {
	return getStatePath(tokenFileName)
}

func getUsernamePath() (string, error) {
	return getStatePath(userFileName)
}

func getStatePath(name string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(home, geniusDirName)
	return filepath.Join(dir, name), nil
}

/* ================= TOKEN ================= */

func GetToken() string {
	if token := strings.TrimSpace(os.Getenv(EnvTokenName)); token != "" {
		return token
	}

	return storedToken()
}

func GetUsername() string {
	if user := strings.TrimSpace(os.Getenv(EnvUserName)); user != "" {
		return user
	}

	return storedUsername()
}

func TokenSource() string {
	if strings.TrimSpace(os.Getenv(EnvTokenName)) != "" {
		return "environment"
	}

	if storedToken() != "" {
		return "file"
	}

	return ""
}

func HasStoredToken() bool {
	return storedToken() != ""
}

func storedToken() string {
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

func storedUsername() string {
	path, err := getUsernamePath()
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
	return SaveAuth(token, "")
}

func SaveAuth(token, username string) error {
	token = strings.TrimSpace(token)
	username = strings.TrimSpace(username)
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
	if err := os.WriteFile(path, []byte(token), 0600); err != nil {
		return err
	}

	if username == "" {
		if userPath, err := getUsernamePath(); err == nil {
			_ = os.Remove(userPath)
		}
		return nil
	}

	userPath, err := getUsernamePath()
	if err != nil {
		return err
	}

	return os.WriteFile(userPath, []byte(username), 0600)
}

func Delete() error {
	path, err := getTokenPath()
	if err != nil {
		return err
	}
	_ = os.Remove(path)

	userPath, err := getUsernamePath()
	if err != nil {
		return nil
	}
	_ = os.Remove(userPath)
	return nil
}
