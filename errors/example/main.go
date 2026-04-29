package main

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidID    = errors.New("invalid user id")
	ErrUserNotFound = errors.New("user not found")
)

type ValidationError struct {
	Field string
	Msg   string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation: %s: %s", e.Field, e.Msg)
}

func main() {
	id := 0

	err := getUserData(id)
	if err != nil {
		fmt.Println("err:", err)

		fmt.Println("Is ErrInvalidID?:", errors.Is(err, ErrInvalidID))
		fmt.Println("Is ErrUserNotFound?:", errors.Is(err, ErrUserNotFound))

		var ve *ValidationError
		fmt.Println("As ValidationError?:", errors.As(err, &ve))
		if ve != nil {
			fmt.Println("field:", ve.Field)
		}
	}
}

func getUserData(id int) error {
	var errs []error
	if id <= 0 {
		errs = append(errs, fmt.Errorf("get user data: %w", ErrInvalidID))
		errs = append(errs, &ValidationError{Field: "id", Msg: "must be positive"})
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	found := false
	if !found {
		return fmt.Errorf("repo: select user id=%d: %w", id, ErrUserNotFound)
	}
	return nil
}

