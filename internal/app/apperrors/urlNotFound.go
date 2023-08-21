// Package apperrors provides custom error definitions for common application-specific scenarios.
package apperrors

import "errors"

// ErrNotFound is an error variable that represents the scenario where the original URL is not found.
var ErrNotFound = errors.New("original URL not found")
