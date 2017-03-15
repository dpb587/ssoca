package authenticated

import (
	"net/http"

	"github.com/dpb587/ssoca/auth"
)

type Requirement struct{}

func (r Requirement) IsSatisfied(_ *http.Request, token *auth.Token) (bool, error) {
	if token == nil {
		return false, nil
	} else if token.ID == "" {
		return false, nil
	}

	return true, nil
}
