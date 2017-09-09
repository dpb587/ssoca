package req

import (
	"net/http"

	"github.com/dpb587/ssoca/auth"
)

//go:generate counterfeiter . RouteHandler
type RouteHandler interface {
	Route() string
	Execute(Request) error
	VerifyAuthorization(*http.Request, *auth.Token) error
}
