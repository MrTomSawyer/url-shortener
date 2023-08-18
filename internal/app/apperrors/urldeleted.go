// Package apperrors provides custom error definitions for common application-specific scenarios.
package apperrors

import "errors"

// ErrURLDeleted is an error variable that represents the scenario where a URL has been deleted.
var ErrURLDeleted = errors.New("the url has been deleted")
