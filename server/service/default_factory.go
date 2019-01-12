package service

import (
	"fmt"

	"github.com/pkg/errors"
)

type defaultFactory struct {
	services map[string]ServiceFactory
}

var _ Factory = defaultFactory{}

func NewDefaultFactory() defaultFactory {
	res := defaultFactory{}
	res.services = map[string]ServiceFactory{}

	return res
}

func (f defaultFactory) Register(serviceFactory ServiceFactory) {
	f.services[serviceFactory.Type()] = serviceFactory
}

func (f defaultFactory) Create(sType string, name string, options map[string]interface{}) (Service, error) {
	factory, ok := f.services[sType]
	if !ok {
		return nil, fmt.Errorf("Unrecognized service type: %s", sType)
	}

	service, err := factory.Create(name, options)
	if err != nil {
		return nil, errors.Wrapf(err, "Creating service %s[%s]", sType, name)
	}

	return service, nil
}
