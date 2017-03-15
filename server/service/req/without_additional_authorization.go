package req

import (
	"net/http"

	"github.com/dpb587/ssoca/auth"
)

type WithoutAdditionalAuthorization struct{}

func (WithoutAdditionalAuthorization) IsAuthorized(r *http.Request, token *auth.Token) (bool, error) {
	return true, nil
}
