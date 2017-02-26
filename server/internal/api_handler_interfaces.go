package internal

import (
	"net/http"

	"github.com/dpb587/ssoca/auth"
)

// @todo rename/suffix these

//go:generate counterfeiter . FakeRouteHandler
type FakeRouteHandler interface {
	Route() string
}

//go:generate counterfeiter . Fake
type Fake interface {
	FakeRouteHandler
	Execute()
}

//go:generate counterfeiter . FakeInHttpResponseWriter
type FakeInHttpResponseWriter interface {
	FakeRouteHandler
	Execute(http.ResponseWriter)
}

//go:generate counterfeiter . FakeInHttpRequest
type FakeInHttpRequest interface {
	FakeRouteHandler
	Execute(*http.Request)
}

//go:generate counterfeiter . FakeInAuthToken
type FakeInAuthToken interface {
	FakeRouteHandler
	Execute(*auth.Token)
}

//go:generate counterfeiter . FakeInApiPayload
type FakeInApiPayload interface {
	FakeRouteHandler
	Execute(map[string]interface{})
}

//go:generate counterfeiter . FakeOutError
type FakeOutError interface {
	FakeRouteHandler
	Execute() error
}

//go:generate counterfeiter . FakeOutInterfaceError
type FakeOutInterfaceError interface {
	FakeRouteHandler
	Execute() (interface{}, error)
}

//go:generate counterfeiter . FakeOutOther
type FakeOutOther interface {
	FakeRouteHandler
	Execute() (string, map[string]interface{}, error)
}
