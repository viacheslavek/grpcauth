package validator

import (
	"errors"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

var (
	ErrEmptyParameter = errors.New("empty parameter")
)

func ValidateEmail(email string) error {
	return validation.Validate(
		email,
		validation.Required,
		is.Email,
	)
}

func ValidatePassword(password string) error {
	return validation.Validate(
		password,
		validation.Required,
		validation.Length(8, 0),
		validation.Match(regexp.MustCompile(`[a-zA-Z]`)).Error("must contain at least one letter"),
		validation.Match(regexp.MustCompile(`[0-9]`)).Error("must contain at least one digit"),
	)
}

func ValidateLogin(login string) error {
	return validation.Validate(
		login,
		validation.Required,
		validation.Match(regexp.MustCompile("^[a-zA-Z0-9]+$")).Error("must be alphanumeric"),
	)
}
