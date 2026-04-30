package main

import (
	"errors"
	"fmt"
)

var (
	ErrMissingName = errors.New("missing name")
	ErrInvalidAge  = errors.New("invalid age")
)

func main() {
	err := validateAll("", -1)

	if err != nil {
		fmt.Println("--- Full Error Message ---")
		fmt.Println(err)

		fmt.Println("\n--- Identifying Individual Errors ---")

		// join 이라고 하더라도 errors.Is는 내부를 순회하며 에러를 찾는다.
		if errors.Is(err, ErrMissingName) {
			fmt.Println("Check: 'Missing Name' error is present.")
		}
		if errors.Is(err, ErrInvalidAge) {
			fmt.Println("Check: 'Invalid Age' error is present.")
		}

		// Go 1.20+ join error -> []error
		fmt.Println("\n--- Iterating Through Error List ---")
		if errs, ok := err.(interface{ Unwrap() []error }); ok {
			for i, e := range errs.Unwrap() {
				fmt.Printf("[%d] Individual Log: %v\n", i+1, e)
			}
		}
	}
}

func validateAll(name string, age int) error {
	var errs []error

	if name == "" {
		errs = append(errs, ErrMissingName)
	}
	if age < 0 {
		errs = append(errs, ErrInvalidAge)
	}

	return errors.Join(errs...)
}
