package service

import (
	"fmt"

	"github.com/dpb587/ssoca/service"
	"github.com/pkg/errors"
)

type defaultFactory struct {
	services map[service.Type]ServiceFactory
}

var _ Factory = defaultFactory{}

func NewDefaultFactory() defaultFactory {
	res := defaultFactory{}
	res.services = map[service.Type]ServiceFactory{}

	return res
}

func (f defaultFactory) Register(serviceFactory ServiceFactory) {
	f.services[serviceFactory.Type()] = serviceFactory
}

func (f defaultFactory) Create(sType service.Type, name string, options map[string]interface{}) (Service, error) {
	factory, ok := f.services[sType]
	if !ok {
		return nil, fmt.Errorf("unrecognized service type: %s", sType)
	}

	service, err := factory.Create(name, options)
	if err != nil {
		return nil, errors.Wrapf(err, "creating service %s[%s]", sType, name)
	}

	return service, nil
}
