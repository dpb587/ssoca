package scope

import (
	"net/http"

	"github.com/dpb587/ssoca/auth"
)

type Requirement struct {
	Present string `yaml:"present"`
}

func (r Requirement) IsSatisfied(_ *http.Request, token auth.Token) (bool, error) {
	if token == nil {
		return false, nil
	} else if !token.HasAttribute(r.Present) {
		return false, nil
	}

	return true, nil
}
