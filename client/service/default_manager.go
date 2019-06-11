package service

import "fmt"

// TODO this is weird; factories are global; services are environment-specific; should be split
type defaultManager struct {
	factories map[string]ServiceFactory
	services  map[string]Service
}

var _ Manager = &defaultManager{}

func NewDefaultManager() Manager {
	res := defaultManager{}
	res.factories = map[string]ServiceFactory{}
	res.services = map[string]Service{}

	return &res
}

func (f *defaultManager) Add(service Service) {
	f.services[service.Name()] = service
}

func (f *defaultManager) AddFactory(factory ServiceFactory) {
	f.factories[factory.Type()] = factory
}

func (f *defaultManager) Get(sType string, sName string) (Service, error) {
	if _, found := f.services[sName]; !found {
		factory, found := f.factories[sType]
		if !found {
			return nil, fmt.Errorf("unrecognized type: %s", sType)
		}

		f.services[sName] = factory.New(sName)
	}

	return f.services[sName], nil
}
