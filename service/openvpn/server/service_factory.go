package server

import (
	"time"

	"github.com/pkg/errors"

	"github.com/dpb587/ssoca/certauth"
	"github.com/dpb587/ssoca/config"
	"github.com/dpb587/ssoca/server/service"
	svc "github.com/dpb587/ssoca/service/openvpn"
	svcconfig "github.com/dpb587/ssoca/service/openvpn/server/config"
)

type ServiceFactory struct {
	svc.ServiceType

	caManager certauth.Manager
}

var _ service.ServiceFactory = ServiceFactory{}

func NewServiceFactory(caManager certauth.Manager) ServiceFactory {
	return ServiceFactory{
		caManager: caManager,
	}
}

func (f ServiceFactory) Create(name string, options map[string]interface{}) (service.Service, error) {
	var cfg svcconfig.Config
	cfg.CertAuth = certauth.NewConfigValue(f.caManager)
	cfg.Validity = 2 * time.Minute

	err := config.RemarshalYAML(options, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "loading config")
	}

	return NewService(name, cfg), nil
}
