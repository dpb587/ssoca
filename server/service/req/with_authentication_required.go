package req

import (
	"net/http"

	"github.com/dpb587/ssoca/auth"
	"github.com/dpb587/ssoca/auth/authz/filter/authenticated"
)

var authz authenticated.Requirement

type WithAuthenticationRequired struct{}

func (WithAuthenticationRequired) VerifyAuthorization(r *http.Request, token *auth.Token) error {
	return authz.VerifyAuthorization(r, token)
}
