package main

import (
	"errors"
	"fmt"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidID    = errors.New("invalid user id")
)

func main() {
	id := 0
	err := getUserData(id)

	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidID):
			fmt.Printf("Handle: bad request (id=%d): %v\n", id, err)
		case errors.Is(err, ErrUserNotFound):
			fmt.Printf("Handle: user not found (id=%d): %v\n", id, err)
		default:
			fmt.Printf("Handle: unknown error: %v\n", err)
		}
		return
	}
	fmt.Println("Success: User data retrieved.")
}

func getUserData(id int) error {
	if id <= 0 {
		return fmt.Errorf("get user data: %w", ErrInvalidID)
	}

	found := false
	if !found {
		return fmt.Errorf("get user data: %w", ErrUserNotFound)
	}
	return nil
}
