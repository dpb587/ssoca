package service

import (
	"fmt"
)

type defaultManager struct {
	services map[string]Service
}

var _ Manager = defaultManager{}

func NewDefaultManager() defaultManager {
	res := defaultManager{}
	res.services = map[string]Service{}

	return res
}

func (f defaultManager) Add(service Service) {
	f.services[service.Name()] = service
}

func (f defaultManager) Get(name string) (Service, error) {
	service, ok := f.services[name]
	if !ok {
		return nil, fmt.Errorf("unrecognized service type: %s", name)
	}

	return service, nil
}

func (f defaultManager) Services() []string {
	services := []string{}

	for name, _ := range f.services {
		services = append(services, name)
	}

	return services
}
