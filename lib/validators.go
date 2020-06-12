package lib

import (
	"errors"
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

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
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("username", username)
	} else {
		panic(errors.New("Error: Could not register validator"))
	}
}
