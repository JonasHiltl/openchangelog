package errs

import "errors"

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

func NewError(domainErr error, appErr error) error {
	return Error{
		appErr:    appErr,
		domainErr: domainErr,
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
