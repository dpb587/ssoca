package or

import (
	"errors"
	"net/http"

	"github.com/dpb587/ssoca/auth"
	"github.com/dpb587/ssoca/auth/authn"
	"github.com/dpb587/ssoca/auth/authz"
	"github.com/dpb587/ssoca/auth/authz/filter"
)

type Requirement struct {
	Requirements []filter.Requirement
}

func (r Requirement) VerifyAuthorization(req *http.Request, token *auth.Token) error {
	for _, requirement := range r.Requirements {
		err := requirement.VerifyAuthorization(req, token)
		if err != nil {
			if _, ok := err.(authn.Error); ok {
				// authentication errors take precedent
				return err
			}
		} else {
			return nil
		}
	}

	return authz.NewError(errors.New("No filters authorized access"))
}
