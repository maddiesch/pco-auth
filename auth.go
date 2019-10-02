package auth

import (
	"context"
	"io/ioutil"
	"log"
)

// PerformInput is the required values for the input of Perform
type PerformInput struct {
	// The Oauth client id from api.planningcenteronline.com
	ClientID string `validate:"required,client_token"`

	// The Oauth client secret from api.planningcenteronline.com
	ClientSecret string `validate:"required,client_token"`

	// The scopes to authenticate with
	Scopes []string `validate:"required,min=1"`

	// The HTTP port to start listening on
	Port int `validate:"required"`

	// An optional place where progress data will be written to.
	Logger *log.Logger

	// Should the ID fetch be skipped? If true the output will be empty strings.
	SkipID bool
}

// PerformOutput is the output of the Perform function call.
type PerformOutput struct {
	// The access token that was just authenticated.
	AccessToken *AccessToken

	// The current person's organization id and person id
	CurrentPerson *CurrentPerson `json:",omitempty"`
}

// Perform will perform the authorization for Planning Center Oauth
func Perform(input *PerformInput) (*PerformOutput, error) {
	return PerformWithContext(context.Background(), input)
}

// PerformWithContext performs the authorization for Planning Center Oauth with the passed context
func PerformWithContext(ctx context.Context, input *PerformInput) (*PerformOutput, error) {
	output := make(chan *PerformOutput)

	newCtx := withContext(ctx)

	go run(newCtx, input, output)

	for {
		select {
		case out := <-output:
			return out, nil
		case <-newCtx.Done():
			return nil, newCtx.Err()
		}
	}
}

func run(ctx context.Context, input *PerformInput, output chan *PerformOutput) {
	if input.Logger == nil {
		input.Logger = newDefaultLogger()
	}

	input.Logger.Printf("Validating input")

	err := Validator().Struct(input)
	if err != nil {
		_fail(ctx, err)
	}

	err = beginAuthorization(input)
	if err != nil {
		_fail(ctx, err)
	}

	callback, err := getCallbackURL(ctx, input.Logger, input.Port)
	if err != nil {
		_fail(ctx, err)
	}

	token, err := AccessTokenFromCallback(&AccessTokenFromCallbackInput{
		ClientID:     input.ClientID,
		ClientSecret: input.ClientSecret,
		CallbackURL:  callback,
		RedirectURL:  redirectURI(input.Port),
		logger:       input.Logger,
	})
	if err != nil {
		_fail(ctx, err)
	}

	var current *CurrentPerson

	if !input.SkipID {
		current, err = FetchCurrentPerson(&FetchCurrentPersonInput{
			AccessToken: token,
			logger:      input.Logger,
		})
		if err != nil {
			_fail(ctx, err)
		}
	}

	input.Logger.Print("Done...")

	output <- &PerformOutput{
		AccessToken:   token,
		CurrentPerson: current,
	}
}

func beginAuthorization(input *PerformInput) error {
	uri, err := AuthorizationURL(&AuthorizationURLInput{
		ClientID:    input.ClientID,
		CallbackURI: redirectURI(input.Port),
		Scopes:      input.Scopes,
	})
	if err != nil {
		return err
	}

	input.Logger.Printf("Opening %s", uri)

	return openBrowser(uri.String())
}

func newDefaultLogger() *log.Logger {
	return log.New(ioutil.Discard, "", 0)
}
