package errors

import (
	"fmt"
	"strings"
)

type ValidationError struct {
	messages []string
}

func (v ValidationError) Error() string {
	return strings.Join(v.messages, "; ")
}

func NewValidationError(field, message string) ValidationError {
	return ValidationError{
		messages: []string{fmt.Sprintf("%s: %s", field, message)},
	}
}

func (v *ValidationError) Add(field, message string) {
	v.messages = append(v.messages, fmt.Sprintf("%s: %s", field, message))
}

func (v ValidationError) HasErrors() bool {
	return len(v.messages) > 0
}
