package service

import (
	"fmt"

	"github.com/pkg/errors"
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
		return nil, fmt.Errorf("Unrecognized service type: %s", name)
	}

	return service, nil
}

func (f defaultManager) GetAuth() (AuthService, error) {
	svc, err := f.Get("auth")
	if err != nil {
		return nil, errors.Wrap(err, "Getting auth service")
	}

	authSvc, ok := svc.(AuthService)
	if !ok {
		return nil, errors.New("Invalid authentication service configured under 'auth' service")
	}

	return authSvc, nil
}

func (f defaultManager) Services() []string {
	services := []string{}

	for name, _ := range f.services {
		services = append(services, name)
	}

	return services
}
