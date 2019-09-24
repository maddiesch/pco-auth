package auth_test

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/maddiesch/pco-auth"
)

func TestAuthorizationURL(t *testing.T) {
	callbackURI, _ := url.Parse("http://foo.bar/callback")

	input := AuthorizationURLInput{
		ClientID:    fakeClientID,
		CallbackURI: callbackURI,
		Scopes:      []string{"people"},
	}

	t.Run("given valid input", func(t *testing.T) {
		url, err := AuthorizationURL(&input)

		require.NoError(t, err)

		expected := fmt.Sprintf("%s://%s/oauth/authorize?%s", PlanningCenterScheme, PlanningCenterHost, `client_id=160547a18450864a2e3f73536e7f76486146acb81b4c4ead3f308621ba044d87&redirect_uri=http%3A%2F%2Ffoo.bar%2Fcallback&response_type=code&scope=people`)

		assert.Equal(t, expected, url.String())
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
