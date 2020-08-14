package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
)

// AccessToken is a Twitch.tv JSON-encoded access token
type AccessToken struct {
	AccessToken  string   `json:"access_token,omitempty"`
	RefreshToken string   `json:"refresh_token,omitempty"`
	ExpiresIn    int      `json:"expires_in,omitempty"`
	Scope        []string `json:"scope,omitempty"`
	TokenType    string   `json:"token_type,omitempty"`
}

// GetAccessToken generates a new access token on Twitch API
func GetAccessToken() (*AccessToken, error) {
	baseURL, err := url.Parse("https://id.twitch.tv/oauth2/token")
	if err != nil {
		return nil, err
	}

	query := url.Values{}

	query.Add("client_id", os.Getenv("CLIENT_ID"))
	query.Add("client_secret", os.Getenv("SECRET_ID"))
	query.Add("grant_type", "client_credentials")
	query.Add("scope", "chat:read chat:edit")

	baseURL.RawQuery = query.Encode()

	response, err := http.Post(baseURL.String(), "", nil)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	accessToken := &AccessToken{}

	err = json.NewDecoder(response.Body).Decode(accessToken)
	if err != nil {
		return nil, err
	}

	return accessToken, nil
}
