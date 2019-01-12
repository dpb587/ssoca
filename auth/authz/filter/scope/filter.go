package scope

import (
	"github.com/pkg/errors"

	"github.com/dpb587/ssoca/auth/authz/filter"
	"github.com/dpb587/ssoca/config"
)

type Filter struct{}

func (f Filter) Create(cfg interface{}) (filter.Requirement, error) {
	requirement := Requirement{}

	err := config.RemarshalYAML(cfg, &requirement)
	if err != nil {
		return nil, errors.Wrap(err, "Loading config")
	}

	if requirement.Present == "" {
		return nil, errors.New("Property must be configured: present")
	}

	return requirement, nil
}
