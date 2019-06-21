package client

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/dpb587/ssoca/client/service"
	globalservice "github.com/dpb587/ssoca/service"
)

const AuthorizationNone = "none"

type AuthTransport struct {
	Runtime Runtime
	Base    http.RoundTripper

	serviceManager service.Manager
}

func NewAuthTransport(runtime Runtime, serviceManager service.Manager, base http.RoundTripper) http.RoundTripper {
	return &AuthTransport{
		Runtime:        runtime,
		Base:           base,
		serviceManager: serviceManager,
	}
}

func (t *AuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Header.Get("Authorization") == AuthorizationNone {
		// do nothing
		req.Header.Del("Authorization")
	} else {
		env, err := t.Runtime.GetEnvironment()
		if err != nil {
			return nil, errors.Wrap(err, "retrieving environment")
		}

		if env.Auth != nil {
			authServiceType := globalservice.Type(env.Auth.Type)

			svc, err := t.serviceManager.Get(authServiceType, env.Auth.Name)
			if err != nil {
				return nil, errors.Wrap(err, "getting authentication service")
			}

			authService, ok := svc.(service.AuthService)
			if !ok {
				return nil, fmt.Errorf("cannot authenticate with service: %s", authServiceType)
			}

			err = authService.AuthRequest(req)
			if err != nil {
				return nil, errors.Wrap(err, "authenticating request")
			}
		}
	}

	return t.base().RoundTrip(req)
}

func (t *AuthTransport) base() http.RoundTripper {
	if t.Base != nil {
		return t.Base
	}
	return http.DefaultTransport
}
