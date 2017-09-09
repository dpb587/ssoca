package remote_ip

import (
	"errors"
	"net"
	"net/http"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/dpb587/ssoca/auth"
	"github.com/dpb587/ssoca/auth/authz"
)

type Requirement struct {
	WithinRaw string    `yaml:"within"`
	Within    net.IPNet `yaml:"-"`
}

func (r Requirement) VerifyAuthorization(req *http.Request, _ *auth.Token) error {
	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return bosherr.WrapError(err, "Parsing remote address")
	}

	ip := net.ParseIP(host)
	if !r.Within.Contains(ip) {
		return authz.NewError(errors.New("Remote IP is not allowed"))
	}

	return nil
}
