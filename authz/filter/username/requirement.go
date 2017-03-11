package username

import (
	"net/http"

	"github.com/dpb587/ssoca/auth"
)

type Requirement struct {
	Is string `yaml:"is"`
}

func (r Requirement) IsSatisfied(_ *http.Request, token *auth.Token) (bool, error) {
	if token == nil {
		return false, nil
	} else if token.Username() == r.Is {
		return true, nil
	}

	return false, nil
}
