package and

import (
	"fmt"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/dpb587/ssoca/authz/filter"
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

	arr, ok := cfg.([]filter.RequireConfig)
	if !ok {
		return nil, fmt.Errorf("Filter options for 'and' is not an array: %#v", cfg)
	}

	for reqIdx, req := range arr {
		if len(req) != 1 {
			return nil, fmt.Errorf("Filter options for item %d of 'and' does not have 1 field", reqIdx)
		}

		for reqType, reqOptions := range req {
			reqFilter, err := f.manager.Get(reqType)
			if err != nil {
				return nil, bosherr.WrapErrorf(err, "Loading filter '%s'", reqType)
			}

			req, err := reqFilter.Create(reqOptions)
			if err != nil {
				return nil, fmt.Errorf("Creating requirement for item %d of 'and'", reqIdx)
			}

			requirement.Requirements = append(requirement.Requirements, req)
		}
	}

	return requirement, nil
}
