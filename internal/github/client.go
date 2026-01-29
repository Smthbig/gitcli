package github

import (
	"errors"
	"net/http"
	"time"
)

const apiBase = "https://api.github.com"

type Client struct {
	http  *http.Client
	token string
	user  string
}

func NewClient() (*Client, error) {
	token := GetToken()
	if token == "" {
		return nil, errors.New("github token not configured")
	}

	user, err := Validate()
	if err != nil {
		return nil, err
	}

	return &Client{
		http:  &http.Client{Timeout: 10 * time.Second},
		token: token,
		user:  user,
	}, nil
}
