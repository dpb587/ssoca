package client

import (
	"fmt"
	"net/http"

	"github.com/dpb587/ssoca/client/service"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
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
			return nil, bosherr.WrapError(err, "Retrieving environment")
		}

		if env.Auth != nil {
			authServiceType := env.Auth.Type

			svc, err := t.serviceManager.Get(authServiceType)
			if err != nil {
				return nil, bosherr.WrapError(err, "Getting authentication service")
			}

			authService, ok := svc.(service.AuthService)
			if !ok {
				return nil, fmt.Errorf("Cannot authenticate with service: %s", authServiceType)
			}

			err = authService.AuthRequest(req)
			if err != nil {
				return nil, bosherr.WrapError(err, "Authenticating request")
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
