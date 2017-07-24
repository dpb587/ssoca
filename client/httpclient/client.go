package httpclient

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	envclient "github.com/dpb587/ssoca/service/env/client"

	"github.com/dpb587/ssoca/client/config"
	"github.com/dpb587/ssoca/client/service"
	baseclient "github.com/dpb587/ssoca/httpclient"
)

type client struct {
	upstream        baseclient.Client
	serviceManager  service.Manager
	configManager   config.Manager
	environmentName string
}

var _ baseclient.Client = client{}

func NewClient(upstream baseclient.Client, serviceManager service.Manager, configManager config.Manager, environmentName string) baseclient.Client {
	return client{
		upstream:        upstream,
		serviceManager:  serviceManager,
		configManager:   configManager,
		environmentName: environmentName,
	}
}

func (c client) APIGet(url string, out interface{}) error {
	err := c.upstream.APIGet(url, out)
	if err != nil && c.attemptReauthenticate(err) == nil {
		return c.upstream.APIGet(url, out)
	}

	return err
}

func (c client) APIPost(url string, in interface{}, out interface{}) error {
	err := c.upstream.APIPost(url, in, out)
	if err != nil && c.attemptReauthenticate(err) == nil {
		return c.upstream.APIPost(url, in, out)
	}

	return err
}

func (c client) Get(url string) (*http.Response, error) {
	res, err := c.upstream.Get(url)
	if err != nil && c.attemptReauthenticate(err) == nil {
		return c.upstream.Get(url)
	}

	return res, err
}

func (c client) Post(url string, contentType string, body io.Reader) (*http.Response, error) {
	res, err := c.upstream.Post(url, contentType, body)
	if err != nil && c.attemptReauthenticate(err) == nil {
		return c.upstream.Post(url, contentType, body)
	}

	return res, err
}

func (c client) attemptReauthenticate(err error) error {
	if !strings.Contains(err.Error(), "HTTP 403") { // TODO lazy improper
		return err
	}

	rawEnvService, err := c.serviceManager.Get("env")
	if err != nil {
		return bosherr.WrapError(err, "Getting env service")
	}

	envService, ok := rawEnvService.(envclient.Service)
	if !ok {
		return bosherr.WrapError(err, "Expecting env service")
	}

	envClient, err := envService.GetClient()
	if err != nil {
		return bosherr.WrapError(err, "Getting env HTTP client")
	}

	envInfo, err := envClient.GetInfo()
	if err != nil {
		return bosherr.WrapError(err, "Getting environment info")
	}

	authServiceType := envInfo.Auth.Type

	svc, err := c.serviceManager.Get(authServiceType)
	if err != nil {
		return bosherr.WrapError(err, "Loading auth service")
	}

	authService, ok := svc.(service.AuthService)
	if !ok {
		return fmt.Errorf("Cannot authenticate with service: %s", authServiceType)
	}

	auth, err := authService.AuthLogin(envInfo.Auth)
	if err != nil {
		return bosherr.WrapError(err, "Authenticating")
	}

	env, err := c.configManager.GetEnvironment(c.environmentName)
	if err != nil {
		return bosherr.WrapError(err, "Getting environment state")
	}

	env.Auth = &config.EnvironmentAuthState{
		Type:    authServiceType,
		Options: auth,
	}

	err = c.configManager.SetEnvironment(env)
	if err != nil {
		return bosherr.WrapError(err, "Updating environment")
	}

	return nil
}
