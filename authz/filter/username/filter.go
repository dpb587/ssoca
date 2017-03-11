package username

import (
	"errors"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/dpb587/ssoca/authz/filter"
	"github.com/dpb587/ssoca/config"
)

type Filter struct{}

func (f Filter) Create(cfg interface{}) (filter.Requirement, error) {
	requirement := Requirement{}

	err := config.RemarshalYAML(cfg, &requirement)
	if err != nil {
		return nil, bosherr.WrapError(err, "Loading config")
	}

	if requirement.Is == "" {
		return nil, errors.New("Property must be configured: is")
	}

	return requirement, nil
}
