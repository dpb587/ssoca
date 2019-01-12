package client

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/dpb587/ssoca/client/service"
	env_api "github.com/dpb587/ssoca/service/env/api"
)

var _ service.AuthService = Service{}

func (s Service) AuthLogin(info env_api.InfoServiceResponse) (interface{}, error) {
	authServiceType := info.Type

	svc, err := s.serviceManager.Get(authServiceType)
	if err != nil {
		return nil, errors.Wrap(err, "Loading auth service")
	}

	authService, ok := svc.(service.AuthService)
	if !ok {
		return nil, fmt.Errorf("Cannot authenticate with service: %s", authServiceType)
	}

	return authService.AuthLogin(info)
}

func (s Service) AuthLogout() error {
	// @todo
	return nil
}

func (s Service) AuthRequest(req *http.Request) error {
	env, err := s.runtime.GetEnvironment()
	if err != nil {
		return errors.Wrap(err, "Getting environment")
	}

	authServiceType := env.Auth.Type

	svc, err := s.serviceManager.Get(authServiceType)
	if err != nil {
		return errors.Wrap(err, "Loading auth service")
	}

	authService, ok := svc.(service.AuthService)
	if !ok {
		return fmt.Errorf("Cannot authenticate with service: %s", authServiceType)
	}

	return authService.AuthRequest(req)
}
