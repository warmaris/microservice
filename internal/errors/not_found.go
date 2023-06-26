package errors

import "fmt"

type NotFoundError struct {
	what string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("not found: %s", e.what)
}

func NewNotFoundError(what string) error {
	return NotFoundError{what: what}
}
