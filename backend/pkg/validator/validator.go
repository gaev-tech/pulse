package validator

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"
)

var _emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// Decode decodes the JSON body of r into destination, disallowing unknown fields.
func Decode(request *http.Request, destination any) error {
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		return errors.New("invalid request body")
	}
	return nil
}

// Email returns an error if the email address is empty or malformed.
func Email(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return errors.New("email is required")
	}
	if !_emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

// Required returns an error if value is empty after trimming whitespace.
func Required(value, field string) error {
	if strings.TrimSpace(value) == "" {
		return errors.New(field + " is required")
	}
	return nil
}
