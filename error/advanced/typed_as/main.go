package main

import (
	"errors"
	"fmt"
)

func main() {
	err := doThing()
	fmt.Println("err:", err)

	var vErr *ValidationError
	// 에러 데이터 추출 (ValidationError -> vErr 할당)
	if errors.As(err, &vErr) {
		fmt.Println("as ValidationError:", "field=", vErr.Field, "msg=", vErr.Msg)
	}

	// 2. 에러 확인
	if errors.Is(err, ErrInvalidInput) {
		fmt.Println("is ErrInvalidInput: true")
	}
}

var ErrInvalidInput = errors.New("invalid input")

type ValidationError struct {
	Field string
	Msg   string
	Err   error
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation: field=%s msg=%s: %v", e.Field, e.Msg, e.Err)
}

func (e *ValidationError) Unwrap() error { return e.Err }

func doThing() error {
	return &ValidationError{
		Field: "email",
		Msg:   "missing @",
		Err:   ErrInvalidInput,
	}
}
