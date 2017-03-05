package server

import (
	"time"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/dpb587/ssoca/certauth"
	"github.com/dpb587/ssoca/config"
	"github.com/dpb587/ssoca/server/service"
	"github.com/dpb587/ssoca/server/service/dynamicvalue"
	svc "github.com/dpb587/ssoca/service/ssh"
	svcconfig "github.com/dpb587/ssoca/service/ssh/config"
)

type ServiceFactory struct {
	caManager           certauth.Manager
	dynamicvalueFactory dynamicvalue.Factory
}

func NewServiceFactory(dynamicvalueFactory dynamicvalue.Factory, caManager certauth.Manager) ServiceFactory {
	return ServiceFactory{
		caManager:           caManager,
		dynamicvalueFactory: dynamicvalueFactory,
	}
}

func (f ServiceFactory) Type() string {
	return svc.Service{}.Type()
}

func (f ServiceFactory) Create(name string, options map[string]interface{}) (service.Service, error) {
	var cfg svcconfig.Config
	cfg.Validity = 2 * time.Minute
	cfg.CertAuth = certauth.NewConfigValue(f.caManager)
	cfg.Principals = dynamicvalue.NewMultiConfigValue(f.dynamicvalueFactory)
	cfg.Target = svcconfig.Target{
		User: dynamicvalue.NewConfigValue(f.dynamicvalueFactory),
	}

	err := config.RemarshalYAML(options, &cfg)
	if err != nil {
		return nil, bosherr.WrapError(err, "Loading config")
	}

	return NewService(name, cfg), nil
}
