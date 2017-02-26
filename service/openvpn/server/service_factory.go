package server

import (
	"strings"
	"time"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/dpb587/ssoca/certauth"
	"github.com/dpb587/ssoca/config"
	"github.com/dpb587/ssoca/server/service"
	svc "github.com/dpb587/ssoca/service/openvpn"
	svcconfig "github.com/dpb587/ssoca/service/openvpn/config"
)

type ServiceFactory struct {
	caManager certauth.Manager
}

func NewServiceFactory(caManager certauth.Manager) ServiceFactory {
	return ServiceFactory{
		caManager: caManager,
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
	config.Profile = strings.TrimSpace(config.Profile)

	ca, err := f.caManager.Get(config.CertAuthName)
	if err != nil {
		return bosherr.WrapError(err, "Getting certificate authority")
	}

	config.CertAuth = ca

	duration, err := time.ParseDuration(config.ValidityString)
	if err != nil {
		return bosherr.WrapError(err, "Parsing validity duration")
	}

	config.Validity = duration

	return nil
}
