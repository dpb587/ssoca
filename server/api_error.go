package server

import (
	"net/http"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type APIError struct {
	error

	Status      int
	PublicError string
}

func NewAPIError(err error, status int, publicError string) APIError {
	apiError, ok := err.(APIError)

	if ok {
		status = apiError.Status
		publicError = apiError.PublicError
	}

	if publicError == "" {
		publicError = http.StatusText(status)
	}

	return APIError{
		error:       err,
		Status:      status,
		PublicError: publicError,
	}
}

func WrapError(err error, msg string) error {
	apiError, ok := err.(APIError)

	if ok {
		return APIError{
			error:       bosherr.WrapError(apiError.error, msg),
			Status:      apiError.Status,
			PublicError: apiError.PublicError,
		}
	}

	return bosherr.WrapError(err, msg)
}

func WrapErrorf(err error, msg string, args ...interface{}) error {
	apiError, ok := err.(APIError)

	if ok {
		return APIError{
			error:       bosherr.WrapErrorf(apiError.error, msg, args...),
			Status:      apiError.Status,
			PublicError: apiError.PublicError,
		}
	}

	return bosherr.WrapErrorf(err, msg, args...)
}
