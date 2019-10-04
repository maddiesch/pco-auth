package auth_test

import (
	"fmt"
	"net/url"
	"testing"

	. "github.com/maddiesch/pco-auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetchCurrentPerson(t *testing.T) {
	t.Run("given a valid input", func(t *testing.T) {
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

		input := &FetchCurrentPersonInput{
			AccessToken: token,
		}

		current, err := FetchCurrentPerson(input)

		require.NoError(t, err)

		assert.Equal(t, "1234", current.OrganizationID)
		assert.Equal(t, "5678", current.PersonID)
	})
}
