package auth_test

import (
	"fmt"
	"net/url"
	"testing"

	. "github.com/maddiesch/pco-auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccessTokenFromCallback(t *testing.T) {
	t.Run("given a valid input", func(t *testing.T) {
		input := &AccessTokenFromCallbackInput{
			CallbackURL: &url.URL{
				Path:     "/callback",
				RawQuery: fmt.Sprintf("code=%s", fakeCode),
			},
			ClientID:     fakeClientID,
			ClientSecret: fakeClientSecret,
			Port:         8080,
		}

		token, err := AccessTokenFromCallback(input)

		require.NoError(t, err)

		assert.Equal(t, "4fe461bac2e2104725d9f2d4f4f0c71d56235e4e089915c18799d3e4e3112e8e", token.Token)
	})
}

func TestPerformRefresh(t *testing.T) {
	t.Run("given a valid refresh_token", func(t *testing.T) {
		token, err := AccessTokenFromCallback(&AccessTokenFromCallbackInput{
			CallbackURL: &url.URL{
				Path:     "/callback",
				RawQuery: fmt.Sprintf("code=%s", fakeCode),
			},
			ClientID:     fakeClientID,
			ClientSecret: fakeClientSecret,
			Port:         8080,
		})

		require.NoError(t, err)

		assert.Equal(t, "4fe461bac2e2104725d9f2d4f4f0c71d56235e4e089915c18799d3e4e3112e8e", token.Token)

		input := &PerformRefreshInput{
			ClientID:     fakeClientID,
			ClientSecret: fakeClientSecret,
		}

		token, err = token.PerformRefresh(input)

		require.NoError(t, err)

		assert.Equal(t, "4f08dfadbbbfd092df3cb9ca4b5c04b1c48a149e15b9866523b1af328a5aa7a7", token.Token)
	})
}
