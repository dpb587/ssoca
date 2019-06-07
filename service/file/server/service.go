package server

import (
	"net/http"

	"github.com/dpb587/ssoca/auth"
	"github.com/dpb587/ssoca/server/service"
	"github.com/dpb587/ssoca/server/service/req"
	svc "github.com/dpb587/ssoca/service/file"
	svcconfig "github.com/dpb587/ssoca/service/file/server/config"
	svcreq "github.com/dpb587/ssoca/service/file/server/req"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type Service struct {
	svc.ServiceType

	name   string
	config svcconfig.Config

	fs boshsys.FileSystem
}

var _ service.Service = Service{}

func NewService(name string, config svcconfig.Config, fs boshsys.FileSystem) Service {
	return Service{
		name:   name,
		config: config,
		fs:     fs,
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
		svcreq.Get{
			Paths: s.config.Paths,
			FS:    s.fs,
		},
		svcreq.List{
			Paths: s.config.Paths,
		},
		svcreq.Metadata{
			Metadata: s.config.Metadata,
		},
	}
}

func (s Service) VerifyAuthorization(_ http.Request, _ *auth.Token) error {
	return nil
}

func (s Service) GetDownloadPaths() []svcconfig.PathConfig {
	return s.config.Paths
}
