package server

import (
	"github.com/pkg/errors"

	svc "github.com/dpb587/ssoca/auth/authn/uaa"
	svcconfig "github.com/dpb587/ssoca/auth/authn/uaa/config"
	"github.com/dpb587/ssoca/config"
	"github.com/dpb587/ssoca/server/service"
)

type ServiceFactory struct{}

func NewServiceFactory() ServiceFactory {
	return ServiceFactory{}
}

func (f ServiceFactory) Type() string {
	return svc.Service{}.Type()
}

func (sf ServiceFactory) Create(name string, options map[string]interface{}) (service.Service, error) {
	var cfg svcconfig.Config

	err := config.RemarshalYAML(options, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "Loading config")
	}

	return NewService(name, cfg), nil
}
