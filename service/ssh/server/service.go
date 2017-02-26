package server

import (
	"net/http"

	"github.com/dpb587/ssoca/auth"
	"github.com/dpb587/ssoca/server/service/req"
	svc "github.com/dpb587/ssoca/service/ssh"
	svcconfig "github.com/dpb587/ssoca/service/ssh/config"
	svcreq "github.com/dpb587/ssoca/service/ssh/server/req"
)

type Service struct {
	svc.Service

	name   string
	config svcconfig.Config
}

func NewService(name string, config svcconfig.Config) Service {
	return Service{
		name:   name,
		config: config,
	}
}

func (s Service) Name() string {
	return s.name
}

func (s Service) Metadata() interface{} {
	return nil
}

func (s Service) GetRoutes() []req.RouteHandler {
	return []req.RouteHandler{
		svcreq.CAPublicKey{
			CertAuth: s.config.CertAuth,
		},
		svcreq.SignPublicKey{
			Validity:        s.config.Validity,
			CriticalOptions: s.config.CriticalOptions,
			Extensions:      s.config.Extensions,
			CertAuth:        s.config.CertAuth,
			Target:          s.config.Target,
		},
	}
}

func (s Service) IsAuthorized(_ http.Request, _ *auth.Token) (bool, error) {
	return true, nil
}
