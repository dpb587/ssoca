package api

import (
	"net/http"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type Error struct {
	error

	Status      int
	PublicError string
}

func NewError(err error, status int, publicError string) Error {
	apiError, ok := err.(Error)

	if ok {
		status = apiError.Status
		publicError = apiError.PublicError
	}

	if publicError == "" {
		publicError = http.StatusText(status)
	}

	return Error{
		error:       err,
		Status:      status,
		PublicError: publicError,
	}
}

func WrapError(err error, msg string) error {
	apiError, ok := err.(Error)

	if ok {
		return Error{
			error:       bosherr.WrapError(apiError.error, msg),
			Status:      apiError.Status,
			PublicError: apiError.PublicError,
		}
	}

	return bosherr.WrapError(err, msg)
}

func WrapErrorf(err error, msg string, args ...interface{}) error {
	apiError, ok := err.(Error)

	if ok {
		return Error{
			error:       bosherr.WrapErrorf(apiError.error, msg, args...),
			Status:      apiError.Status,
			PublicError: apiError.PublicError,
		}
	}

	return bosherr.WrapErrorf(err, msg, args...)
}
