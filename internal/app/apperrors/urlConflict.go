// Package apperrors provides custom error definitions for common application-specific scenarios.
package apperrors

import (
	"fmt"
)

// URLConflict is a custom error type representing a conflict scenario with an HTTP 409 status.
// It includes a Value field to hold related data and an Err field to embed a lower-level error.
type URLConflict struct {
	Value string // Value is the related data associated with the conflict.
	Err   error  // Err is an embedded error that provides additional context.
}

// NewURLConflict is a constructor function that creates a new URLConflict error instance.
func NewURLConflict(value string, err error) error {
	return &URLConflict{Value: value, Err: err}
}

// Error implements the error interface for the URLConflict type.
func (uc *URLConflict) Error() string {
	return fmt.Sprintf("conflict: an original URL already exists. Shortened URL: %s", uc.Value)
}
