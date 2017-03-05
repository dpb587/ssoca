package dynamicvalue

import (
	"net/http"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/dpb587/ssoca/auth"
)

type MultiAnyValue []Value

func (mav MultiAnyValue) Evaluate(arg0 *http.Request, arg1 *auth.Token) ([]string, error) {
	values := []string{}

	for _, value := range mav {
		res, err := value.Evaluate(arg0, arg1)
		if err != nil {
			return nil, bosherr.WrapError(err, "Evaluating template")
		}

		values = append(values, res)
	}

	return values, nil
}
