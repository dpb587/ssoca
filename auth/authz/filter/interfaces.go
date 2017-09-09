package filter

import (
	"net/http"

	"github.com/dpb587/ssoca/auth"
)

//go:generate counterfeiter . Requirement
type Requirement interface {
	VerifyAuthorization(*http.Request, *auth.Token) error
}

//go:generate counterfeiter . Manager
type Manager interface {
	Add(string, Filter)
	Get(string) (Filter, error)
	Filters() []string
}

//go:generate counterfeiter . Filter
type Filter interface {
	Create(interface{}) (Requirement, error)
}
