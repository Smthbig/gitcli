package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

/* ================= CHECK ================= */

func RepoExists(owner, repo string) (bool, error) {
	c, err := NewClient()
	if err != nil {
		return false, err
	}

	url := fmt.Sprintf("%s/repos/%s/%s", apiBase, owner, repo)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "token "+c.token)

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

	return false, fmt.Errorf("github api error: %s", resp.Status)
}

/* ================= CREATE ================= */

func CreateRepo(owner, repo string, private bool) error {
	c, err := NewClient()
	if err != nil {
		return err
	}

	payload := map[string]any{
		"name":    repo,
		"private": private,
	}

	body, _ := json.Marshal(payload)

	var url string
	if owner == c.user {
		url = apiBase + "/user/repos"
	} else {
		url = fmt.Sprintf("%s/orgs/%s/repos", apiBase, owner)
	}

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "token "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		return fmt.Errorf("failed to create repo: %s", resp.Status)
	}

	return nil
}
