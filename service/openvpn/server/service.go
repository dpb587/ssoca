package server

import (
	"net/http"

	"github.com/dpb587/ssoca/auth"
	"github.com/dpb587/ssoca/server/service"
	"github.com/dpb587/ssoca/server/service/req"
	svc "github.com/dpb587/ssoca/service/openvpn"
	svcconfig "github.com/dpb587/ssoca/service/openvpn/config"
	svcreq "github.com/dpb587/ssoca/service/openvpn/server/req"
)

type Service struct {
	svc.Service

	name   string
	config svcconfig.Config
}

var _ service.Service = Service{}

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
		svcreq.SignUserCSR{
			Validity:    s.config.Validity,
			CertAuth:    s.config.CertAuth,
			BaseProfile: s.config.Profile,
		},
		svcreq.BaseProfile{
			BaseProfile: s.config.Profile,
			CertAuth:    s.config.CertAuth,
		},
	}
}

func (s Service) IsAuthorized(_ http.Request, _ *auth.Token) (bool, error) {
	return true, nil
}
