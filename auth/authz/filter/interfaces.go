package filter

import (
	"net/http"

	"github.com/dpb587/ssoca/auth"
)

//go:generate counterfeiter . Requirement
type Requirement interface {
	IsSatisfied(*http.Request, *auth.Token) (bool, error)
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
