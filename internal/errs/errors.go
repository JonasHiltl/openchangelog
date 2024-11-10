package errs

import (
	"errors"
	"net/http"
)

// Generic errors that can be wrapped in the domain
var (
	ErrBadRequest         = errors.New("bad request")
	ErrNotFound           = errors.New("not found")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrServiceUnavailable = errors.New("service unavailable")
)

type Error struct {
	// The specific app error or the service
	appErr error
	// The generic error type
	domainErr error
}

func (e Error) Status() int {
	switch e.domainErr {
	case ErrBadRequest:
		return http.StatusBadRequest
	case ErrNotFound:
		return http.StatusNotFound
	case ErrUnauthorized:
		return http.StatusUnauthorized
	case ErrServiceUnavailable:
		return http.StatusServiceUnavailable
	}

	return http.StatusInternalServerError
}

func (e Error) Msg() string {
	return e.appErr.Error()
}

func NewError(domainErr error, appErr error) error {
	return Error{
		appErr:    appErr,
		domainErr: domainErr,
	}
}

func NewBadRequest(wrapped error) error {
	return Error{
		appErr:    wrapped,
		domainErr: ErrBadRequest,
	}
}

func NewNotFound(wrapped error) error {
	return Error{
		appErr:    wrapped,
		domainErr: ErrNotFound,
	}
}

func NewUnauthorized(wrapped error) error {
	return Error{
		appErr:    wrapped,
		domainErr: ErrUnauthorized,
	}
}

func NewServiceUnavailable(wrapped error) error {
	return Error{
		appErr:    wrapped,
		domainErr: ErrServiceUnavailable,
	}
}

func (e Error) AppErr() error {
	return e.appErr
}

func (e Error) DomainErr() error {
	return e.domainErr
}

func (e Error) Error() string {
	return errors.Join(e.domainErr, e.appErr).Error()
}
