package rest

import (
	"errors"
	"net/http"

	"github.com/jonashiltl/openchangelog/internal/domain"
)

// Error represents a handler error. It provides methods for a HTTP status
// code and embeds the built-in error interface.
type Error interface {
	error
	Status() int
}

// RestError represents an error with an associated HTTP status code.
type RestError struct {
	Code int
	Err  error
}

// Allows StatusError to satisfy the error interface.
func (se RestError) Error() string {
	return se.Err.Error()
}

// Returns our HTTP status code.
func (se RestError) Status() int {
	return se.Code
}

func RestErrorFromDomain(err error) RestError {
	statErr := RestError{
		Code: http.StatusInternalServerError,
		Err:  err,
	}

	var domErr domain.Error
	if errors.As(err, &domErr) {
		statErr.Err = domErr.AppErr()
		switch domErr.DomainErr() {
		case domain.ErrBadRequest:
			statErr.Code = http.StatusBadRequest
		case domain.ErrNotFound:
			statErr.Code = http.StatusNotFound
		}
	}

	return statErr
}
