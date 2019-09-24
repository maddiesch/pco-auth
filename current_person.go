package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// CurrentPerson represents the authenticated person
type CurrentPerson struct {
	// The authenticated person's organization ID
	OrganizationID string

	// The authenticated person's ID
	PersonID string
}

// FetchCurrentPersonInput contains the required parameters for a call to FetchCurrentPerson
type FetchCurrentPersonInput struct {
	AccessToken *AccessToken `validate:"required"`
	logger      *log.Logger
}

// FetchCurrentPerson performs a GET to the Planning Center API to get the current person information.
func FetchCurrentPerson(input *FetchCurrentPersonInput) (*CurrentPerson, error) {
	if input.logger == nil {
		input.logger = newDefaultLogger()
	}

	err := Validator().Struct(input)
	if err != nil {
		return nil, err
	}

	input.logger.Print("Fetching OrganizationID and PersonID")

	req, err := http.NewRequest("GET", apiURL("/check").String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", input.AccessToken.Token))

	body, err := sendRequest(req, input.logger)
	if err != nil {
		return nil, err
	}

	output := struct {
		Data struct {
			Attributes struct {
				PersonID       string `json:"person_id"`
				OrganizationID string `json:"organization_id"`
			} `json:"attributes"`
		} `json:"data"`
	}{}

	err = json.Unmarshal(body, &output)
	if err != nil {
		return nil, err
	}

	return &CurrentPerson{
		PersonID:       output.Data.Attributes.PersonID,
		OrganizationID: output.Data.Attributes.OrganizationID,
	}, nil
}
