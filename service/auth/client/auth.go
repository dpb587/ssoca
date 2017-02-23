package client

import (
	"fmt"
	"net/http"

	"github.com/dpb587/ssoca/client/service"
	env_api "github.com/dpb587/ssoca/service/env/api"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

func (s Service) AuthLogin(info env_api.InfoServiceResponse) (interface{}, error) {
	authServiceType := info.Type

	svc, err := s.serviceManager.Get(authServiceType)
	if err != nil {
		return nil, bosherr.WrapError(err, "Loading auth service")
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
		return bosherr.WrapError(err, "Getting environment")
	}

	authServiceType := env.Auth.Type

	svc, err := s.serviceManager.Get(authServiceType)
	if err != nil {
		return bosherr.WrapError(err, "Loading auth service")
	}

	authService, ok := svc.(service.AuthService)
	if !ok {
		return fmt.Errorf("Cannot authenticate with service: %s", authServiceType)
	}

	return authService.AuthRequest(req)
}
