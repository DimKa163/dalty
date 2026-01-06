package protoerr

import (
	"fmt"
	"strings"
)

type ErrorType int

const (
	Internal ErrorType = iota
	InvalidArgs
)

func (err ErrorType) String() string {
	return [...]string{"Internal", "InvalidArgs"}[err]
}

type ServiceError struct {
	Type    ErrorType
	Message string
	err     error
}

func mewInternalError(err error) error {
	return &ServiceError{
		Type:    Internal,
		Message: "internal error occurred",
		err:     err,
	}
}

func newInvalidArgsError(message string, args []string) error {
	return &ServiceError{
		Type:    InvalidArgs,
		Message: fmt.Sprintf("%s. invalid arguments(%s)", message, strings.Join(args, ", ")),
	}
}
func (e *ServiceError) Error() string {
	return fmt.Sprintf("%s %s", e.Type.String(), e.Message)
}
