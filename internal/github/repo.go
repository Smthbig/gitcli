package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

/*
CreateRepo:
- Works for user OR org
- Lazily fetches authenticated user
- No RepoExists required
- Proper error parsing
*/

func CreateRepo(owner, repo string, private bool) error {
	c, err := NewClient()
	if err != nil {
		return err
	}

	authUser, err := c.GetAuthenticatedUser()
	if err != nil {
		return err
	}

	payload := map[string]any{
		"name":    repo,
		"private": private,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	var url string
	if owner == authUser {
		url = apiBase + "/user/repos"
	} else {
		url = fmt.Sprintf("%s/orgs/%s/repos", apiBase, owner)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "git-genius")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 201 {
		return nil
	}

	// Read body for detailed error
	bodyBytes, _ := io.ReadAll(resp.Body)
	msg := string(bodyBytes)

	switch resp.StatusCode {
	case 401:
		return fmt.Errorf("invalid or expired GitHub token")

	case 403:
		return fmt.Errorf("permission denied for owner '%s'", owner)

	case 422:
		return fmt.Errorf("repository may already exist or validation failed")

	default:
		return fmt.Errorf("github api error (%d): %s",
			resp.StatusCode,
			msg,
		)
	}
}
func RepoExists(owner, repo string) (bool, error) {

	c, err := NewClient()
	if err != nil {
		return false, err
	}

	url := fmt.Sprintf("%s/repos/%s/%s", apiBase, owner, repo)

	req, err := c.newRequest("GET", url, nil)
	if err != nil {
		return false, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return true, nil
	}

	if resp.StatusCode == 404 {
		return false, nil
	}

	return false, fmt.Errorf("github api error (status %d)", resp.StatusCode)
}
