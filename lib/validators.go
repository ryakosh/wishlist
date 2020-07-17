package lib

import (
	"errors"
	"regexp"

	"github.com/go-playground/validator/v10"
)

var ErrValidationFailed = errors.New("Field validation failed")

var Validator = validator.New()

var (
	rgxUsername = regexp.MustCompile("^[a-z0-9_-]+$")
)

func username(fl validator.FieldLevel) bool {
	if uname, ok := fl.Field().Interface().(string); ok {
		if rgxUsername.MatchString(uname) {
			return true
		}
	}
	return false
}

func init() {
	Validator.RegisterValidation("username", username)
}
