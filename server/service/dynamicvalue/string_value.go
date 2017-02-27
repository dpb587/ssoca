package dynamicvalue

import (
	"net/http"

	"github.com/dpb587/ssoca/auth"
)

type stringValue string

func NewStringValue(val string) stringValue {
	return stringValue(val)
}

func (cv stringValue) Evaluate(_ *http.Request, _ *auth.Token) (string, error) {
	return string(cv), nil
}
