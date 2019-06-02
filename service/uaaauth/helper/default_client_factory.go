package helper

import (
	boshuaa "github.com/cloudfoundry/bosh-cli/uaa"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"github.com/pkg/errors"
)

type DefaultClientFactory struct{}

func (cm DefaultClientFactory) CreateClient(url string, authClient string, authClientSecret string, caCertificate string) (boshuaa.UAA, error) {
	config, err := boshuaa.NewConfigFromURL(url)
	if err != nil {
		return nil, errors.Wrap(err, "parsing UAA URL")
	}

	config.CACert = caCertificate
	config.Client = authClient
	config.ClientSecret = authClientSecret

	factory := boshuaa.NewFactory(boshlog.NewLogger(boshlog.LevelNone))
	client, err := factory.New(config)
	if err != nil {
		return nil, errors.Wrap(err, "creating UAA client")
	}

	return client, nil
}
