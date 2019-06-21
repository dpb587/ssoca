package server

import (
	"crypto/x509"
	"encoding/pem"
	"strings"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"

	"github.com/dpb587/ssoca/certauth"
	"github.com/dpb587/ssoca/config"
	"github.com/dpb587/ssoca/server/service"
	"github.com/dpb587/ssoca/server/service/dynamicvalue"
	svc "github.com/dpb587/ssoca/service/ssh"
	svcconfig "github.com/dpb587/ssoca/service/ssh/server/config"
)

type ServiceFactory struct {
	svc.ServiceType

	caManager           certauth.Manager
	dynamicvalueFactory dynamicvalue.Factory
}

var _ service.ServiceFactory = ServiceFactory{}

func NewServiceFactory(dynamicvalueFactory dynamicvalue.Factory, caManager certauth.Manager) ServiceFactory {
	return ServiceFactory{
		caManager:           caManager,
		dynamicvalueFactory: dynamicvalueFactory,
	}
}

func (f ServiceFactory) Create(name string, options map[string]interface{}) (service.Service, error) {
	var cfg svcconfig.Config
	cfg.Validity = 2 * time.Minute
	cfg.CertAuth = certauth.NewConfigValue(f.caManager)
	cfg.Extensions = svcconfig.ExtensionDefaults
	cfg.Principals = dynamicvalue.NewMultiConfigValue(f.dynamicvalueFactory)
	cfg.CriticalOptions = svcconfig.NewCriticalOptions(f.dynamicvalueFactory)
	cfg.Target = svcconfig.Target{
		User: dynamicvalue.NewConfigValue(f.dynamicvalueFactory),
	}

	err := config.RemarshalYAML(options, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "loading config")
	}

	if cfg.Target.PublicKey != "" && strings.Contains(cfg.Target.PublicKey, "-----") {
		publicKeyPEM, _ := pem.Decode([]byte(cfg.Target.PublicKey))
		if publicKeyPEM == nil {
			return nil, errors.New("failed to parse public key")
		}

		rsa, err := x509.ParsePKIXPublicKey(publicKeyPEM.Bytes)
		if err != nil {
			return nil, errors.Wrap(err, "parsing public key")
		}

		publicKey, err := ssh.NewPublicKey(rsa)
		if err != nil {
			return nil, errors.Wrap(err, "parsing ssh public key")
		}

		cfg.Target.PublicKey = string(ssh.MarshalAuthorizedKey(publicKey))
	}

	return NewService(name, cfg), nil
}
