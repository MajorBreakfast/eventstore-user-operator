package eventstore

import "errors"

var (
	// ErrUserNotFound is returned when the user is not found
	ErrUserNotFound = errors.New("User was not found")
)
