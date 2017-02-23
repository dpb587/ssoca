package remote_ip

import (
	"net"
	"net/http"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/dpb587/ssoca/auth"
)

type Requirement struct {
	WithinRaw string    `yaml:"within"`
	Within    net.IPNet `yaml:"-"`
}

func (r Requirement) IsSatisfied(req *http.Request, _ auth.Token) (bool, error) {
	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return false, bosherr.WrapError(err, "Parsing remote address")
	}

	ip := net.ParseIP(host)
	if !r.Within.Contains(ip) {
		return false, nil
	}

	return true, nil
}
