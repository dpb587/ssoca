package and

import (
	"net/http"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/dpb587/ssoca/auth"
	"github.com/dpb587/ssoca/authz/filter"
)

type Requirement struct {
	Requirements []filter.Requirement
}

func (r Requirement) IsSatisfied(req *http.Request, token *auth.Token) (bool, error) {
	for _, requirement := range r.Requirements {
		satisfied, err := requirement.IsSatisfied(req, token)
		if err != nil {
			return false, bosherr.WrapError(err, "Evaluating requirements")
		} else if satisfied != true {
			return false, nil
		}
	}

	return true, nil
}
