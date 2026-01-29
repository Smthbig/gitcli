package github

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"git-genius/internal/system"
)

const (
	geniusDir = ".git/.genius"
	tokenFile = geniusDir + "/token"
	apiUser   = "https://api.github.com/user"
)

type userResponse struct {
	Login string `json:"login"`
}

/* ================= TOKEN ================= */

func GetToken() string {
	data, _ := os.ReadFile(tokenFile)
	return string(data)
}

func Save(token string) error {
	if token == "" {
		return errors.New("empty token")
	}
	_ = os.MkdirAll(geniusDir, 0700)
	return os.WriteFile(tokenFile, []byte(token), 0600)
}

func Delete() {
	_ = os.Remove(tokenFile)
}

/* ================= VALIDATION ================= */

func Validate() (string, error) {
	token := GetToken()
	if token == "" {
		return "", errors.New("no github token")
	}

	if !system.Online {
		return "offline-mode", nil
	}

	client := http.Client{Timeout: 5 * time.Second}
	req, _ := http.NewRequest("GET", apiUser, nil)
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("User-Agent", "git-genius")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.New("invalid github token")
	}

	var u userResponse
	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return "", err
	}

	return u.Login, nil
}
