package auth_test

import (
	"log"
	"os"
	"testing"

	. "github.com/maddiesch/pco-auth"
)

var (
	TestCredentials FakeConfigOutput
)

func TestMain(m *testing.M) {
	TestCredentials = SetupForTesting(&FakeConfig{
		OrganizationID: "1234",
		PersonID:       "5678",
		Logger:         log.New(os.Stderr, "[fake server] -> ", 0),
	})

	os.Exit(m.Run())
}
