package server

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	svc "github.com/dpb587/ssoca/authn/uaa"
	svcconfig "github.com/dpb587/ssoca/authn/uaa/config"
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
		return nil, bosherr.WrapError(err, "Loading config")
	}

	return NewService(name, cfg), nil
}
