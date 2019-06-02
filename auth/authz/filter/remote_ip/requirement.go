package remote_ip

import (
	"net"
	"net/http"

	"github.com/pkg/errors"

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
		return errors.Wrap(err, "parsing remote address")
	}

	ip := net.ParseIP(host)
	if !r.Within.Contains(ip) {
		return authz.NewError(errors.New("remote IP is not allowed"))
	}

	return nil
}
