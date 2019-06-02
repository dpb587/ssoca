package remote_ip

import (
	"fmt"
	"net"
	"strings"

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

	if requirement.WithinRaw == "" {
		return nil, errors.New("property must be configured: within")
	}

	netmask := requirement.WithinRaw

	if !strings.Contains(netmask, "/") {
		ip := net.ParseIP(netmask)

		if ip.To4() != nil {
			netmask = fmt.Sprintf("%s/32", netmask)
		} else {
			netmask = fmt.Sprintf("%s/128", netmask)
		}
	}

	_, cidr, err := net.ParseCIDR(netmask)
	if err != nil {
		return nil, errors.Wrap(err, "parsing CIDR")
	}

	requirement.Within = *cidr

	return requirement, nil
}
