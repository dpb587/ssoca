package client

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/dpb587/ssoca/client/service"
)

var _ service.AuthService = Service{}

func (s Service) AuthLogin() error {
	// TODO fetch default service and execute AuthLogin()
	return errors.New("TODO")
}

func (s Service) AuthLogout() error {
	// @todo
	return nil
}

func (s Service) AuthRequest(req *http.Request) error {
	env, err := s.runtime.GetEnvironment()
	if err != nil {
		return errors.Wrap(err, "getting environment")
	}

	authServiceType := env.Auth.Type

	svc, err := s.serviceManager.Get(authServiceType)
	if err != nil {
		return errors.Wrap(err, "loading auth service")
	}

	authService, ok := svc.(service.AuthService)
	if !ok {
		return fmt.Errorf("cannot authenticate with service: %s", authServiceType)
	}

	return authService.AuthRequest(req)
}
