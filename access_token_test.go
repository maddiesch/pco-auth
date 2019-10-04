package auth_test

import (
	"fmt"
	"net/url"
	"strings"
	"testing"

	. "github.com/maddiesch/pco-auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertHavePrefix(t *testing.T, prefix string, value string) {
	if strings.HasPrefix(value, prefix) {
		return
	}

	t.Logf("Expected '%s' to have the prefix '%s'", value, prefix)
	t.Fail()
}

func TestAccessTokenFromCallback(t *testing.T) {
	t.Run("given a valid input", func(t *testing.T) {
		input := &AccessTokenFromCallbackInput{
			CallbackURL: &url.URL{
				Path:     "/callback",
				RawQuery: fmt.Sprintf("code=%s", TestCredentials.Code),
			},
			ClientID:     TestCredentials.ClientID,
			ClientSecret: TestCredentials.ClientSecret,
			RedirectURL:  &url.URL{},
		}

		token, err := AccessTokenFromCallback(input)

		require.NoError(t, err)

		assertHavePrefix(t, "fake_0", token.Token)
	})
}

func TestPerformRefresh(t *testing.T) {
	t.Run("given a valid refresh_token", func(t *testing.T) {
		token, err := AccessTokenFromCallback(&AccessTokenFromCallbackInput{
			CallbackURL: &url.URL{
				Path:     "/callback",
				RawQuery: fmt.Sprintf("code=%s", TestCredentials.Code),
			},
			ClientID:     TestCredentials.ClientID,
			ClientSecret: TestCredentials.ClientSecret,
			RedirectURL:  &url.URL{},
		})

		require.NoError(t, err)

		assertHavePrefix(t, "fake_0", token.Token)

		before := token.Token

		input := &PerformRefreshInput{
			ClientID:     TestCredentials.ClientID,
			ClientSecret: TestCredentials.ClientSecret,
		}

		token, err = token.PerformRefresh(input)

		require.NoError(t, err)

		assertHavePrefix(t, "fake_0", token.Token)

		assert.NotEqual(t, before, token.Token)
	})
}
