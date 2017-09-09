package username

import (
	"errors"
	"net/http"

	"github.com/dpb587/ssoca/auth"
	"github.com/dpb587/ssoca/auth/authn"
	"github.com/dpb587/ssoca/auth/authz"
)

type Requirement struct {
	Is string `yaml:"is"`
}

func (r Requirement) VerifyAuthorization(_ *http.Request, token *auth.Token) error {
	if token == nil {
		return authn.NewError(errors.New("Authentication token missing"))
	} else if token.Username() == r.Is {
		return nil
	}

	return authz.NewError(errors.New("Username does not match"))
}
