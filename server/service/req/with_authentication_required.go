package req

import (
	"net/http"

	"github.com/dpb587/ssoca/auth"
	"github.com/dpb587/ssoca/auth/authz/filter/authenticated"
)

var authz authenticated.Requirement

type WithAuthenticationRequired struct{}

func (WithAuthenticationRequired) IsAuthorized(r *http.Request, token *auth.Token) (bool, error) {
	return authz.IsSatisfied(r, token)
}
