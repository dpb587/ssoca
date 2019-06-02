package service

import "fmt"

type defaultManager struct {
	services map[string]Service
}

var _ Manager = defaultManager{}

func NewDefaultManager() Manager {
	res := defaultManager{}
	res.services = map[string]Service{}

	return res
}

func (f defaultManager) Add(service Service) {
	f.services[service.Type()] = service
}

func (f defaultManager) Get(sType string) (Service, error) {
	service, ok := f.services[sType]
	if !ok {
		return nil, fmt.Errorf("unrecognized type: %s", sType)
	}

	return service, nil
}

func (f defaultManager) Services() []string {
	services := []string{}

	for sType, _ := range f.services {
		services = append(services, sType)
	}

	return services
}
