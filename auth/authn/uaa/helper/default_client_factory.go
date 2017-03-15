package helper

import (
	boshuaa "github.com/cloudfoundry/bosh-cli/uaa"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

type DefaultClientFactory struct{}

func (cm DefaultClientFactory) CreateClient(url string, authClient string, authClientSecret string, caCertificate string) (boshuaa.UAA, error) {
	config, err := boshuaa.NewConfigFromURL(url)
	if err != nil {
		return nil, bosherr.WrapError(err, "Parsing UAA URL")
	}

	config.CACert = caCertificate
	config.Client = authClient
	config.ClientSecret = authClientSecret

	factory := boshuaa.NewFactory(boshlog.NewLogger(boshlog.LevelNone))
	client, err := factory.New(config)
	if err != nil {
		return nil, bosherr.WrapError(err, "Creating UAA client")
	}

	return client, nil
}
