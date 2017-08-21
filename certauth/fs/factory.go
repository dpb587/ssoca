package fs

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/sirupsen/logrus"

	"github.com/dpb587/ssoca/certauth"
	"github.com/dpb587/ssoca/config"
)

type Factory struct {
	fs     boshsys.FileSystem
	logger logrus.FieldLogger
}

var _ certauth.ProviderFactory = Factory{}

func NewFactory(fs boshsys.FileSystem, logger logrus.FieldLogger) Factory {
	return Factory{
		fs:     fs,
		logger: logger,
	}
}

func (f Factory) Create(name string, options map[string]interface{}) (certauth.Provider, error) {
	var cfg Config

	err := config.RemarshalYAML(options, &cfg)
	if err != nil {
		return nil, bosherr.WrapError(err, "Loading config")
	}

	provider := NewProvider(
		name,
		cfg,
		f.fs,
		f.logger.WithFields(logrus.Fields{
			"certauth.name": name,
		}),
	)

	return provider, nil
}
