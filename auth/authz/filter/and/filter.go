package and

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/dpb587/ssoca/auth/authz/filter"
	"github.com/dpb587/ssoca/config"
)

type Filter struct {
	manager filter.Manager
}

func NewFilter(manager filter.Manager) Filter {
	return Filter{
		manager: manager,
	}
}

func (f Filter) Create(cfg interface{}) (filter.Requirement, error) {
	requirement := Requirement{}

	arr := []filter.RequireConfig{}

	err := config.RemarshalYAML(cfg, &arr)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse 'and' config")
	}

	for reqIdx, req := range arr {
		if len(req) != 1 {
			return nil, fmt.Errorf("Filter options for item %d of 'and' does not have 1 field", reqIdx)
		}

		for reqType, reqOptions := range req {
			reqFilter, err := f.manager.Get(reqType)
			if err != nil {
				return nil, errors.Wrapf(err, "Loading filter '%s'", reqType)
			}

			req, err := reqFilter.Create(reqOptions)
			if err != nil {
				return nil, errors.Wrapf(err, "Creating requirement for item %d of 'and'", reqIdx)
			}

			requirement.Requirements = append(requirement.Requirements, req)
		}
	}

	return requirement, nil
}
