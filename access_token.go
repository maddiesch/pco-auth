package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

// AccessToken is the Oauth access token used for authorization of requests
type AccessToken struct {
	Token     string
	Kind      string
	ExpiresIn int64
	Refresh   string
	Scope     string
	CreatedAt int64
}

// SignRequest adds the appropriate Authorization header to the request.
func (a *AccessToken) SignRequest(req *http.Request) {
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", a.Kind, a.Token))
}

// PerformRefreshInput contains the required information from the refresh token
type PerformRefreshInput struct {
	// The Oauth client id from api.planningcenteronline.com
	ClientID string `validate:"required,client_token"`

	// The Oauth client secret from api.planningcenteronline.com
	ClientSecret string `validate:"required,client_token"`

	logger *log.Logger
}

// PerformRefresh refreshes the access token and returns a new one.
func (a *AccessToken) PerformRefresh(input *PerformRefreshInput) (*AccessToken, error) {
	if input.logger == nil {
		input.logger = newDefaultLogger()
	}

	err := Validator().Struct(input)
	if err != nil {
		return nil, err
	}

	input.logger.Printf("Performing token refresh")

	body, err := json.Marshal(map[string]string{
		"grant_type":    "refresh_token",
		"client_id":     input.ClientID,
		"client_secret": input.ClientSecret,
		"refresh_token": a.Refresh,
	})
	if err != nil {
		return nil, err
	}

	return sendAccessTokenRequest(body, input.logger)
}

type inboundAccessToken struct {
	Token     string `json:"access_token"`
	Kind      string `json:"token_type"`
	ExpiresIn int64  `json:"expires_in"`
	Refresh   string `json:"refresh_token"`
	Scope     string `json:"scope"`
	CreatedAt int64  `json:"created_at"`
}

// AccessTokenFromCallbackInput is the required parameters for a call to AccessTokenFromCallback
type AccessTokenFromCallbackInput struct {
	CallbackURL  *url.URL `validate:"required"`
	ClientID     string   `validate:"required"`
	ClientSecret string   `validate:"required"`
	RedirectURL  *url.URL `validate:"required"`
	logger       *log.Logger
}

// AccessTokenFromCallback returns a new access token using the callback URL
func AccessTokenFromCallback(input *AccessTokenFromCallbackInput) (*AccessToken, error) {
	if input.logger == nil {
		input.logger = newDefaultLogger()
	}

	err := Validator().Struct(input)
	if err != nil {
		return nil, err
	}

	if input.CallbackURL.Query().Get("code") == "" {
		return nil, errors.New("callback url missing 'code' parameter")
	}

	input.logger.Printf("Fetching access token")

	body, err := json.Marshal(map[string]string{
		"grant_type":    "authorization_code",
		"code":          input.CallbackURL.Query().Get("code"),
		"client_id":     input.ClientID,
		"client_secret": input.ClientSecret,
		"redirect_uri":  input.RedirectURL.String(),
	})
	if err != nil {
		return nil, err
	}

	return sendAccessTokenRequest(body, input.logger)
}

func sendAccessTokenRequest(body []byte, logger *log.Logger) (*AccessToken, error) {
	uri := apiURL("/oauth/token")

	req, err := http.NewRequest("POST", uri.String(), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	body, err = sendRequest(req, logger)
	if err != nil {
		return nil, err
	}

	it := inboundAccessToken{}
	err = json.Unmarshal(body, &it)
	if err != nil {
		return nil, err
	}

	token := AccessToken(it)

	return &token, nil
}
