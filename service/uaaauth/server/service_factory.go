package server

import (
	"github.com/pkg/errors"

	"github.com/dpb587/ssoca/config"
	"github.com/dpb587/ssoca/server/service"
	svc "github.com/dpb587/ssoca/service/uaaauth"
	svcconfig "github.com/dpb587/ssoca/service/uaaauth/config"
)

type ServiceFactory struct {
	svc.ServiceType
}

func NewServiceFactory() ServiceFactory {
	return ServiceFactory{}
}

func (sf ServiceFactory) Create(name string, options map[string]interface{}) (service.Service, error) {
	var cfg svcconfig.Config

	err := config.RemarshalYAML(options, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "loading config")
	}

	return NewService(name, cfg), nil
}
