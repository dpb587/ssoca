package certauth

import (
	"fmt"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type defaultFactory struct {
	providers map[string]ProviderFactory
}

func NewDefaultFactory() defaultFactory {
	res := defaultFactory{}
	res.providers = map[string]ProviderFactory{}

	return res
}

func (f defaultFactory) Register(name string, providerFactory ProviderFactory) {
	f.providers[name] = providerFactory
}

func (f defaultFactory) Create(name string, sType string, options map[string]interface{}) (Provider, error) {
	factory, ok := f.providers[sType]
	if !ok {
		return nil, fmt.Errorf("Unrecognized provider type: %s", sType)
	}

	provider, err := factory.Create(name, options)
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Creating provider %s", sType)
	}

	return provider, nil
}
