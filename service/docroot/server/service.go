package server

import (
	"fmt"
	"net/http"

	"github.com/dpb587/ssoca/auth"
	"github.com/dpb587/ssoca/server/service"
	"github.com/dpb587/ssoca/server/service/req"
	svc "github.com/dpb587/ssoca/service/docroot"
	svcconfig "github.com/dpb587/ssoca/service/docroot/server/config"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type Service struct {
	svc.ServiceType

	name   string
	config svcconfig.Config

	fs boshsys.FileSystem
}

var _ service.Service = Service{}

func NewService(name string, config svcconfig.Config, fs boshsys.FileSystem) *Service {
	return &Service{
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
		req.RouteHandlerFunc{
			Func: http.StripPrefix(fmt.Sprintf("/%s/", s.Name()), http.FileServer(http.Dir(s.config.AbsPath))).ServeHTTP,
		},
	}
}

func (s Service) VerifyAuthorization(_ http.Request, _ *auth.Token) error {
	return nil
}
