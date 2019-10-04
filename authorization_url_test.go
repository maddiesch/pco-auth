package auth_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/maddiesch/pco-auth"
)

func TestAuthorizationURL(t *testing.T) {
	callbackURI, _ := url.Parse("http://foo.bar/callback")

	input := AuthorizationURLInput{
		ClientID:    TestCredentials.ClientID,
		CallbackURI: callbackURI,
		Scopes:      []string{"people"},
	}

	t.Run("given valid input", func(t *testing.T) {
		uri, err := AuthorizationURL(&input)

		require.NoError(t, err)

		query := url.Values{}
		query.Set("client_id", TestCredentials.ClientID)
		query.Set("redirect_uri", "http://foo.bar/callback")
		query.Set("response_type", "code")
		query.Set("scope", "people")

		expected := &url.URL{
			Scheme:   PlanningCenterScheme,
			Host:     PlanningCenterHost,
			Path:     "/oauth/authorize",
			RawQuery: query.Encode(),
		}

		assert.Equal(t, expected.String(), uri.String())
	})

	t.Run("given invalid input for client id", func(t *testing.T) {
		bad := AuthorizationURLInput(input)
		bad.ClientID = "foo-bar"
		_, err := AuthorizationURL(&bad)

		require.Error(t, err)

		assert.Equal(t, "Key: 'AuthorizationURLInput.ClientID' Error:Field validation for 'ClientID' failed on the 'client_token' tag", err.Error())
	})

	t.Run("given invalid callback uri", func(t *testing.T) {
		bad := AuthorizationURLInput(input)
		bad.CallbackURI = nil
		_, err := AuthorizationURL(&bad)

		require.Error(t, err)

		assert.Equal(t, "Key: 'AuthorizationURLInput.CallbackURI' Error:Field validation for 'CallbackURI' failed on the 'required' tag", err.Error())
	})

	t.Run("given no scopes", func(t *testing.T) {
		bad := AuthorizationURLInput(input)
		bad.Scopes = []string{}
		_, err := AuthorizationURL(&bad)

		require.Error(t, err)

		assert.Equal(t, "Key: 'AuthorizationURLInput.Scopes' Error:Field validation for 'Scopes' failed on the 'min' tag", err.Error())
	})
}
