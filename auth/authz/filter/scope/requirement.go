package scope

import (
	"errors"
	"net/http"

	"github.com/dpb587/ssoca/auth"
	"github.com/dpb587/ssoca/auth/authn"
	"github.com/dpb587/ssoca/auth/authz"
)

type Requirement struct {
	Present string `yaml:"present"`
}

func (r Requirement) VerifyAuthorization(_ *http.Request, token *auth.Token) error {
	if token == nil {
		return authn.NewError(errors.New("authentication token missing"))
	}

	for _, scope := range token.Groups {
		if scope != r.Present {
			continue
		}

		return nil
	}

	return authz.NewError(errors.New("scope is missing"))
}
