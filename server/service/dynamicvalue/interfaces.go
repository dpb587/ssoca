package dynamicvalue

import (
	"net/http"

	"github.com/dpb587/ssoca/auth"
)

//go:generate counterfeiter . Factory
type Factory interface {
	Create(string) (Value, error)
}

//go:generate counterfeiter . Value
type Value interface {
	Evaluate(*http.Request, *auth.Token) (string, error)
}

//go:generate counterfeiter . MultiValue
type MultiValue interface {
	Evaluate(*http.Request, *auth.Token) ([]string, error)
}
