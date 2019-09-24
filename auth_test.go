package auth_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	. "github.com/maddiesch/pco-auth"
)

const (
	fakeClientID     = "160547a18450864a2e3f73536e7f76486146acb81b4c4ead3f308621ba044d87"
	fakeClientSecret = "bad04c191cc60e910500564359658445e6e2c8122fcaea2807222771c235bc89"
	fakeCode         = "a9119e83815085054d1112a3f8343f89abc294a53f630b603ace0234ea5ebe8b"
)

var (
	fakeAuthenticateTokenResponse = []byte(`{"access_token":"4fe461bac2e2104725d9f2d4f4f0c71d56235e4e089915c18799d3e4e3112e8e","token_type":"Bearer","expires_in":7200,"refresh_token":"57be4d5957038a1d29ef8c0f6481b3e0e5cc9ae0757348edb61aeae8034d836c","scope":"people services","created_at":1569361077}`)
	fakeRefreshTokenResponse      = []byte(`{"access_token":"4f08dfadbbbfd092df3cb9ca4b5c04b1c48a149e15b9866523b1af328a5aa7a7","token_type":"Bearer","expires_in":7200,"refresh_token":"8a6cfbe93eb9e54674a901ef0e0ab168367a181e7893e786b6289ee88b2fcf24","scope":"people services","created_at":1569361077}`)
	fakeCheckResponse             = []byte(`{"data":{"id":"1234","type":"AuthenticationCheck","attributes":{"person_id":"5678","organization_id":"1234","auth_method":"OAuth2","scopes":["people","services"],"oauth2_application_name":"CLI Testing"}}}`)
)

func TestMain(m *testing.M) {
	fake := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth/token":
			if r.Method != "POST" {
				w.WriteHeader(http.StatusMethodNotAllowed)
			} else {
				body, _ := ioutil.ReadAll(r.Body)
				defer r.Body.Close()

				input := struct {
					GrantType string `json:"grant_type"`
				}{}

				_ = json.Unmarshal(body, &input)

				if input.GrantType == "authorization_code" {
					w.Write(fakeAuthenticateTokenResponse)
				} else if input.GrantType == "refresh_token" {
					w.Write(fakeRefreshTokenResponse)
				} else {
					w.WriteHeader(http.StatusBadRequest)
				}
			}
		case "/check":
			if r.Method != "GET" {
				w.WriteHeader(http.StatusMethodNotAllowed)
			} else if r.Header.Get("Authorization") != "Bearer 4fe461bac2e2104725d9f2d4f4f0c71d56235e4e089915c18799d3e4e3112e8e" {
				w.WriteHeader(http.StatusUnauthorized)
			} else {
				w.Write(fakeCheckResponse)
			}
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))

	fakeURL, _ := url.Parse(fake.URL)

	HTTPClient = fake.Client()
	PlanningCenterScheme = fakeURL.Scheme
	PlanningCenterHost = fakeURL.Host

	os.Exit(m.Run())
}
