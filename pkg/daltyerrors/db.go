package daltyerrors

import "errors"

var ErrNotFound = errors.New("not found")

type StorageError struct {
	Message string
	Value   []any
	Err     error
}

func (e *StorageError) Error() string {
	return e.Message
}

func (e *StorageError) Unwrap() error {
	return e.Err
}

func (e *StorageError) Is(target error) bool {
	return errors.Is(e.Unwrap(), target)
}

func NewNotFoundError(err error, message string, value ...any) error {
	return &StorageError{
		Message: message,
		Value:   value,
		Err:     err,
	}
}
