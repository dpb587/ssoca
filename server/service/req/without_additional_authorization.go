package req

import (
	"net/http"

	"github.com/dpb587/ssoca/auth"
)

type WithoutAdditionalAuthorization struct{}

func (WithoutAdditionalAuthorization) VerifyAuthorization(r *http.Request, token *auth.Token) error {
	return nil
}
