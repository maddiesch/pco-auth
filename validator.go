package auth

import (
	"regexp"
	"strings"
	"sync"

	"gopkg.in/go-playground/validator.v9"
)

var (
	validatorSetup    sync.Once
	validatorInstance *validator.Validate
	tokenRegex        = regexp.MustCompile(`\A[a-fA-F0-9]{64}\z`)
)

// Validator returns the shared validator instance
func Validator() *validator.Validate {
	validatorSetup.Do(func() {
		validatorInstance = validator.New()
		validatorInstance.RegisterValidation("client_token", func(fl validator.FieldLevel) bool {
			if strings.HasPrefix(fl.Field().String(), "fake_0") { // Fake testing token
				return true
			}
			return tokenRegex.MatchString(fl.Field().String())
		})
	})
	return validatorInstance
}
