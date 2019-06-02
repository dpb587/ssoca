package certauth

import (
	"fmt"

	"github.com/pkg/errors"
)

type defaultFactory struct {
	providers map[string]ProviderFactory
}

var _ Factory = defaultFactory{}

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
		return nil, fmt.Errorf("unrecognized provider type: %s", sType)
	}

	provider, err := factory.Create(name, options)
	if err != nil {
		return nil, errors.Wrapf(err, "creating provider %s", sType)
	}

	return provider, nil
}
