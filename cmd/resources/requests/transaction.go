package requests

import (
	"errors"
	"strings"
	"unicode/utf8"
)

type TransactionRequest struct {
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}

// Validate performs basic validation checks on the MovieRequest struct.
func (tr *TransactionRequest) Validate() error {
	var required []string

	if utf8.RuneCountInString(tr.Description) > 1000 {
		return errors.New("field 'description' should be 1000 characters max")
	}
	if tr.Description == "" {
		required = append(required, "description")
	}

	if tr.Amount <= 0 {
		return errors.New("field 'amount' should be a positive number greater than 0")
	}

	if len(required) > 0 {
		return errors.New(strings.Join(required, ", ") + " required.")
	}
	return nil
}
