package apperrors

import "fmt"

type URLConflict struct {
	Value string
	Err   error
}

func (uc *URLConflict) Error() string {
	return fmt.Sprintf("conflict: such original URL already exists. Shortened URL: %s", uc.Value)
}

func NewURLConflict(value string, err error) error {
	return &URLConflict{value, err}
}
