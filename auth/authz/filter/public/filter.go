package public

import "github.com/dpb587/ssoca/auth/authz/filter"

type Filter struct{}

func (f Filter) Create(cfg interface{}) (filter.Requirement, error) {
	return Requirement{}, nil
}
