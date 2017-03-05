package dynamicvalue

import (
	"net/http"

	"github.com/dpb587/ssoca/auth"
)

type stringValue string

func NewStringValue(val string) Value {
	return stringValue(val)
}

func (cv stringValue) Evaluate(_ *http.Request, _ *auth.Token) (string, error) {
	return string(cv), nil
}
