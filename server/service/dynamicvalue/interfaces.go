package dynamicvalue

import (
	"net/http"

	"github.com/dpb587/ssoca/auth"
)

type Factory interface {
	Create(string) (Value, error)
}

type Value interface {
	Evaluate(*http.Request, *auth.Token) (string, error)
}
