package main

import (
	"errors"
	"fmt"
)

var ErrUserNotFound = errors.New("user not found")

func main() {
	id := 42

	err := getUserData(id)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			fmt.Printf("Handle: user not found (id=%d): %v\n", id, err)
			return
		}
		fmt.Printf("Handle: unknown error: %v\n", err)
		return
	}

	fmt.Println("Success")
}

func getUserData(id int) error {
	// 예시: DB에서 조회했는데 없음
	found := false
	if !found {
		return fmt.Errorf("repo: select user id=%d: %w", id, ErrUserNotFound)
	}
	return nil
}

