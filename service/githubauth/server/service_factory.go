package server

import (
	"github.com/dpb587/ssoca/config"
	"github.com/dpb587/ssoca/server/service"
	svc "github.com/dpb587/ssoca/service/githubauth"
	svcconfig "github.com/dpb587/ssoca/service/githubauth/config"
	"github.com/pkg/errors"
)

type ServiceFactory struct {
	svc.ServiceType

	endpointURL string
}

func NewServiceFactory(endpointURL string) ServiceFactory {
	return ServiceFactory{
		endpointURL: endpointURL,
	}
}

func (f ServiceFactory) Create(name string, options map[string]interface{}) (service.Service, error) {
	var cfg svcconfig.Config

	err := config.RemarshalYAML(options, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "loading config")
	}

	cfg.AbsolutizeRedirects(f.endpointURL)

	return NewService(name, f.endpointURL, cfg), nil
}
