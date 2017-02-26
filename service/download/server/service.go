package server

import (
	"net/http"

	"github.com/dpb587/ssoca/auth"
	"github.com/dpb587/ssoca/server/service/req"
	svc "github.com/dpb587/ssoca/service/download"
	svcconfig "github.com/dpb587/ssoca/service/download/config"
	svcreq "github.com/dpb587/ssoca/service/download/server/req"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type Service struct {
	svc.Service

	name   string
	config svcconfig.Config

	fs boshsys.FileSystem
}

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
		svcreq.List{
			Paths: s.config.Paths,
		},
		svcreq.Get{
			Paths: s.config.Paths,
			FS:    s.fs,
		},
	}
}

func (s Service) IsAuthorized(_ http.Request, _ *auth.Token) (bool, error) {
	return true, nil
}

func (s Service) GetDownloadPaths() []svcconfig.PathConfig {
	return s.config.Paths
}
