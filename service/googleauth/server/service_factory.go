package server

import (
	"github.com/pkg/errors"

	"github.com/dpb587/ssoca/config"
	serverconfig "github.com/dpb587/ssoca/server/config"
	"github.com/dpb587/ssoca/server/service"
	svc "github.com/dpb587/ssoca/service/googleauth"
	svcconfig "github.com/dpb587/ssoca/service/googleauth/server/config"
)

type ServiceFactory struct {
	svc.ServiceType

	endpointURL string
	redirects   serverconfig.ServerRedirectConfig
}

func NewServiceFactory(endpointURL string, redirects serverconfig.ServerRedirectConfig) ServiceFactory {
	return ServiceFactory{
		endpointURL: endpointURL,
		redirects:   redirects,
	}
}

func (f ServiceFactory) Create(name string, options map[string]interface{}) (service.Service, error) {
	var cfg svcconfig.Config

	err := config.RemarshalYAML(options, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "loading config")
	}

	cfg.ApplyRedirectDefaults(f.redirects.AuthSuccess, f.redirects.AuthFailure)

	return NewService(name, f.endpointURL, cfg), nil
}
