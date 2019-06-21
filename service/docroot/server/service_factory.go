package server

import (
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/pkg/errors"

	"github.com/dpb587/ssoca/config"
	"github.com/dpb587/ssoca/server/service"
	svc "github.com/dpb587/ssoca/service/docroot"
	svcconfig "github.com/dpb587/ssoca/service/docroot/server/config"
)

type ServiceFactory struct {
	svc.ServiceType

	fs boshsys.FileSystem
}

var _ service.ServiceFactory = ServiceFactory{}

func NewServiceFactory(fs boshsys.FileSystem) ServiceFactory {
	return ServiceFactory{
		fs: fs,
	}
}

func (f ServiceFactory) Create(name string, options map[string]interface{}) (service.Service, error) {
	var cfg svcconfig.Config

	err := config.RemarshalYAML(options, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "loading config")
	}

	err = f.validateConfig(&cfg)
	if err != nil {
		return nil, errors.Wrap(err, "validating config")
	}

	return NewService(name, cfg, f.fs), nil
}

func (f ServiceFactory) validateConfig(config *svcconfig.Config) error {
	absPath, err := f.fs.ExpandPath(config.Path)
	if err != nil {
		return errors.Wrap(err, "expanding path")
	}

	config.AbsPath = absPath

	return nil
}
