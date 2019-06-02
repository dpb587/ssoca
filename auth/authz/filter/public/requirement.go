package public

import (
	"net/http"

	"github.com/dpb587/ssoca/auth"
)

type Requirement struct{}

func (r Requirement) VerifyAuthorization(_ *http.Request, token *auth.Token) error {
	return nil
}
