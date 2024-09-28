package resources

import (
	"errors"
	"net/mail"
	"strings"
	"unicode/utf8"
)

// User represents a user of this application
type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

// Validate performs basic validation checks on the User struct.
func (u *User) Validate() error {
	var required []string

	if u.Name == "" {
		required = append(required, "name")
	}
	if u.Email == "" {
		required = append(required, "email")
	}

	if utf8.RuneCountInString(u.Name) > 50 {
		return errors.New("field 'name' should be 50 characters max")
	}
	_, err := mail.ParseAddress(u.Email)
	if err != nil {
		return errors.New("invalid email format")

	}

	if u.Password == "" {
		required = append(required, "password")
	}

	if len(required) > 0 {
		return errors.New(strings.Join(required, ", ") + " required.")
	}
	return nil
}
