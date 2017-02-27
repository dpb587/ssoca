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

	err := config.RemarshalYAML(options, &cfg)
	if err != nil {
		return nil, bosherr.WrapError(err, "Loading config")
	}

	err = f.validateConfig(&cfg)
	if err != nil {
		return nil, bosherr.WrapError(err, "Validating config")
	}

	return NewService(name, cfg), nil
}

func (f ServiceFactory) validateConfig(config *svcconfig.Config) error {
	ca, err := f.caManager.Get(config.CertAuthName)
	if err != nil {
		return bosherr.WrapError(err, "Getting certificate authority")
	}

	config.CertAuth = ca

	duration, err := time.ParseDuration(config.ValidityString)
	if err != nil {
		return bosherr.WrapError(err, "Parsing duration")
	}

	config.Validity = duration

	if config.Target.RawUser != "" {
		value, err := f.dynamicvalueFactory.Create(config.Target.RawUser)
		if err != nil {
			return bosherr.WrapError(err, "Parsing dynamic value target.user")
		}

		config.Target.User = value
	}

	return nil
}
