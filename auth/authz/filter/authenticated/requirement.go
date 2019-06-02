package authenticated

import (
	"errors"
	"net/http"

	"github.com/dpb587/ssoca/auth"
	"github.com/dpb587/ssoca/auth/authn"
)

type Requirement struct{}

func (r Requirement) VerifyAuthorization(_ *http.Request, token *auth.Token) error {
	if token == nil {
		return authn.NewError(errors.New("authentication token missing"))
	} else if token.ID == "" {
		return authn.NewError(errors.New("authentication ID missing"))
	}

	return nil
}
