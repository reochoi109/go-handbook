package main

import (
	"errors"
	"strings"
)

func Add(a, b int) int {
	return a + b
}

func Div(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("divide by zero")
	}
	return a / b, nil
}

func NormalizeName(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}

func IsEven(n int) bool {
	return n%2 == 0
}
