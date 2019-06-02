package cli

import (
	"fmt"

	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"

	clientcmd "github.com/dpb587/ssoca/client/cmd"
	"github.com/dpb587/ssoca/client/config"
	"github.com/dpb587/ssoca/client/service"
	envclient "github.com/dpb587/ssoca/service/env/client"
)

type Login struct {
	*clientcmd.ServiceCommand `no-flag:"true"`

	SkipVerify bool `long:"skip-verify" description:"Skip verification of authentication, once complete"`

	ServiceManager service.Manager
	GetClient      GetClient
}

var _ flags.Commander = Login{}

func (c Login) Execute(_ []string) error {
	rawEnvService, err := c.ServiceManager.Get("env")
	if err != nil {
		return errors.Wrap(err, "getting env service")
	}

	envService, ok := rawEnvService.(envclient.Service)
	if !ok {
		return errors.Wrap(err, "expecting env service")
	}

	envClient, err := envService.GetClient()
	if err != nil {
		return errors.Wrap(err, "getting env HTTP client")
	}

	envInfo, err := envClient.GetInfo()
	if err != nil {
		return errors.Wrap(err, "getting environment info")
	}

	authServiceType := envInfo.Auth.Type

	svc, err := c.ServiceManager.Get(authServiceType)
	if err != nil {
		return errors.Wrap(err, "loading auth service")
	}

	authService, ok := svc.(service.AuthService)
	if !ok {
		return fmt.Errorf("cannot authenticate with service: %s", authServiceType)
	}

	auth, err := authService.AuthLogin(envInfo.Auth)
	if err != nil {
		return errors.Wrap(err, "authenticating")
	}

	env, err := c.Runtime.GetEnvironment()
	if err != nil {
		return errors.Wrap(err, "getting environment state")
	}

	env.Auth = &config.EnvironmentAuthState{
		Type:    authServiceType,
		Options: auth,
	}

	configManager, err := c.Runtime.GetConfigManager()
	if err != nil {
		return errors.Wrap(err, "getting config manager")
	}

	err = configManager.SetEnvironment(env)
	if err != nil {
		return errors.Wrap(err, "updating environment")
	}

	if c.SkipVerify {
		return nil
	}

	err = c.verify()
	if err != nil {
		return errors.Wrap(err, "verifying authentication")
	}

	return nil
}

func (c Login) verify() error {
	ui := c.Runtime.GetUI()

	client, err := c.GetClient()
	if err != nil {
		return errors.Wrap(err, "getting client")
	}

	authInfo, err := client.GetInfo()
	if err != nil {
		return errors.Wrap(err, "getting remote authentication info")
	}

	if authInfo.ID == "" {
		return errors.New("failed to use authentication credentials")
	}

	ui.PrintBlock([]byte(fmt.Sprintf("Successfully logged in as %s\n", authInfo.ID)))

	return nil
}
