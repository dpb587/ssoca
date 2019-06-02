package username

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
		return nil, errors.Wrap(err, "loading config")
	}

	if requirement.Is == "" {
		return nil, errors.New("property must be configured: is")
	}

	return requirement, nil
}
