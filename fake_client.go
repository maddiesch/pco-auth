package auth

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"time"
)

// FakeConfig is the configuration of the fake server
type FakeConfig struct {
	Logger         *log.Logger
	Scopes         []string
	PersonID       string
	OrganizationID string
}

// FakeConfigOutput contains the values needed to perform successful requests to the fake server
type FakeConfigOutput struct {
	ClientID     string
	ClientSecret string
	Code         string
	URL          *url.URL
}

// SetupForTesting performs the setup for pco-auth to stub all requests to Planning Center
// and return fake data.
//
// You don't haven't to do anything with the returned URL but it's there if you want to use
// this as a fake server for something else
func SetupForTesting(config *FakeConfig) FakeConfigOutput {
	if config.PersonID == "" {
		config.PersonID = "-1"
	}
	if config.OrganizationID == "" {
		config.OrganizationID = "-1"
	}
	if config.Scopes == nil || len(config.Scopes) == 0 {
		config.Scopes = []string{"people"}
	}
	if config.Logger == nil {
		config.Logger = log.New(ioutil.Discard, "", 0)
	}

	generate := func() string {
		buffer := make([]byte, 29)
		rand.Read(buffer)
		return fmt.Sprintf("fake_0%x", buffer)
	}

	var accessToken string
	var refreshToken string

	output := FakeConfigOutput{
		ClientID:     generate(),
		ClientSecret: generate(),
		Code:         generate(),
	}

	oauthTokenHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			config.Logger.Println("The HTTP method is not a POST")
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		body, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		input := struct {
			GrantType    string `json:"grant_type"`
			ClientID     string `json:"client_id"`
			ClientSecret string `json:"client_secret"`
			Code         string `json:"code"`
			RefreshToken string `json:"refresh_token"`
		}{}

		_ = json.Unmarshal(body, &input)

		if input.ClientID != output.ClientID {
			config.Logger.Println("The request's client_id does not match the expected client_id")
			w.WriteHeader(http.StatusBadRequest)
		}

		if input.ClientSecret != output.ClientSecret {
			config.Logger.Println("The request's client_secret does not match the expected client_secret")
			w.WriteHeader(http.StatusBadRequest)
		}

		if input.GrantType == "refresh_token" {
			if input.RefreshToken != refreshToken {
				config.Logger.Println("The request's refresh token does not match the expected value")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		accessToken = generate()
		refreshToken = generate()
		body, _ = json.Marshal(map[string]interface{}{
			"access_token":  accessToken,
			"token_type":    "Bearer",
			"expires_in":    7200,
			"refresh_token": refreshToken,
			"scope":         strings.Join(config.Scopes, " "),
			"created_at":    time.Now().Unix(),
		})
		w.Write(body)
	}

	checkHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			config.Logger.Println("The HTTP method is not a GET")
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if accessToken == "" || r.Header.Get("Authorization") != fmt.Sprintf("Bearer %s", accessToken) {
			config.Logger.Println("The Authorization header does not match the expected value")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		data, _ := json.Marshal(map[string]interface{}{
			"data": map[string]interface{}{
				"id":   config.OrganizationID,
				"type": "AuthenticationCheck",
				"attributes": map[string]interface{}{
					"person_id":               config.PersonID,
					"organization_id":         config.OrganizationID,
					"auth_method":             "Oauth2",
					"scopes":                  config.Scopes,
					"oauth2_application_name": "pco-auth testing server",
				},
			},
		})

		w.Write(data)
	}

	fake := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/oauth/token":
			oauthTokenHandler(w, r)
		case "/check":
			checkHandler(w, r)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))

	output.URL, _ = url.Parse(fake.URL)

	HTTPClient = fake.Client()
	PlanningCenterScheme = output.URL.Scheme
	PlanningCenterHost = output.URL.Host

	return output
}
