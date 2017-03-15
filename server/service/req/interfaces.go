package req

import (
	"net/http"

	"github.com/dpb587/ssoca/auth"
)

//go:generate counterfeiter . RouteHandler
type RouteHandler interface {
	Route() string
	Execute(Request) error
	IsAuthorized(*http.Request, *auth.Token) (bool, error)
}
