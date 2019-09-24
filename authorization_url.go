package auth

import (
	"net/url"
	"strings"
)

// AuthorizationURLInput contains the required parameters for building the authorization URL
type AuthorizationURLInput struct {
	ClientID    string   `validate:"required,client_token"`
	CallbackURI *url.URL `validate:"required"`
	Scopes      []string `validate:"min=1"`
}

// AuthorizationURL returns the URL to open in a web browser to start authorization
func AuthorizationURL(input *AuthorizationURLInput) (*url.URL, error) {
	err := Validator().Struct(input)
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	query.Set("client_id", input.ClientID)
	query.Set("redirect_uri", input.CallbackURI.String())
	query.Set("response_type", "code")
	query.Set("scope", strings.Join(input.Scopes, " "))

	uri := apiURL("/oauth/authorize")
	uri.RawQuery = query.Encode()

	return uri, nil
}
