package main

import (
	"errors"
	"fmt"
)

var ErrInvalidInput = errors.New("invalid input")

type ValidationError struct {
	Field string
	Msg   string
	Err   error
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("field %s: %s", e.Field, e.Msg)
}

func (e *ValidationError) Unwrap() error { return e.Err }

func RegisterUser(email string) error {
	if email == "" {
		return &ValidationError{
			Field: "email",
			Msg:   "Email address is a required field.",
			Err:   ErrInvalidInput,
		}
	}
	return nil
}

func main() {
	err := RegisterUser("")

	if err != nil {
		var vErr *ValidationError
		if errors.As(err, &vErr) {
			fmt.Printf("HTTP 400 - Client Response: [%s] Field Error: %s\n", vErr.Field, vErr.Msg)
		}

		if errors.Is(err, ErrInvalidInput) {
			fmt.Println("System Log: Request rejected due to invalid input values.")
		}
	}
}
