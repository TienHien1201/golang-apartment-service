package xgoogle

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

type Profile struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

type Client struct {
	clientID     string
	clientSecret string
	callBackURL  string
}

func New(clientID, secret, callback string) *Client {
	return &Client{clientID, secret, callback}
}

func (c *Client) AuthURL() string {
	u := url.Values{}
	u.Set("client_id", c.clientID)
	u.Set("redirect_uri", c.callBackURL)
	u.Set("response_type", "code")
	u.Set("scope", "openid profile email")
	u.Set("access_type", "offline")
	u.Set("prompt", "consent")

	return "https://accounts.google.com/o/oauth2/v2/auth?" + u.Encode()
}

func (c *Client) GetProfile(ctx context.Context, accessToken string) (*Profile, error) {
	req, _ := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://www.googleapis.com/oauth2/v2/userinfo",
		nil,
	)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var profile Profile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, err
	}

	return &profile, nil
}

func (c *Client) ExchangeCode(ctx context.Context, code string) (string, error) {
	data := url.Values{}
	data.Set("client_id", c.clientID)
	data.Set("client_secret", c.clientSecret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", c.callBackURL)

	req, _ := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://oauth2.googleapis.com/token",
		strings.NewReader(data.Encode()),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.AccessToken, nil
}
