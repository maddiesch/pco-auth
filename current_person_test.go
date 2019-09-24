package auth_test

import (
	"testing"

	. "github.com/maddiesch/pco-auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetchCurrentPerson(t *testing.T) {
	t.Run("given a valid input", func(t *testing.T) {
		input := &FetchCurrentPersonInput{
			AccessToken: &AccessToken{
				Token: "4fe461bac2e2104725d9f2d4f4f0c71d56235e4e089915c18799d3e4e3112e8e",
			},
		}

		current, err := FetchCurrentPerson(input)

		require.NoError(t, err)

		assert.Equal(t, "1234", current.OrganizationID)
		assert.Equal(t, "5678", current.PersonID)
	})
}
