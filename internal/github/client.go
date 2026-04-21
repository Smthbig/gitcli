package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const apiBase = "https://api.github.com"

type Client struct {
	http  *http.Client
	token string
	user  string // lazily resolved
}

/*
NewClient:
- Lightweight
- No validation here
- Safe for repeated calls
*/
func NewClient() (*Client, error) {
	token := GetToken()
	return NewClientWithToken(token)
}

func NewClientWithToken(token string) (*Client, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, errors.New("github token not configured")
	}

	return &Client{
		http: &http.Client{
			Timeout: 12 * time.Second,
		},
		token: token,
	}, nil
}

func ValidateToken(token string) (string, error) {
	client, err := NewClientWithToken(token)
	if err != nil {
		return "", err
	}
	return client.GetAuthenticatedUser()
}

/*
newRequest:
- Centralized request builder
- Ensures headers are always correct
*/
func (c *Client) newRequest(method, url string, body io.Reader) (*http.Request, error) {

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("User-Agent", "git-genius")
	req.Header.Set("Accept", "application/vnd.github+json")

	return req, nil
}

/*
getAuthenticatedUser:
- Lazy resolution
- Cached after first call
*/
func (c *Client) GetAuthenticatedUser() (string, error) {

	if c.user != "" {
		return c.user, nil
	}

	req, err := c.newRequest("GET", apiBase+"/user", nil)
	if err != nil {
		return "", err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf(
			"invalid or expired GitHub token (status %d)",
			resp.StatusCode,
		)
	}

	var data struct {
		Login string `json:"login"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	c.user = data.Login
	return c.user, nil
}
