package cmd

import (
	"errors"
	"fmt"

	"github.com/dpb587/ssoca/client/config"
	"github.com/dpb587/ssoca/client/service"
	"github.com/jessevdk/go-flags"

	clientcmd "github.com/dpb587/ssoca/client/cmd"
	envclient "github.com/dpb587/ssoca/service/env/client"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type Login struct {
	clientcmd.ServiceCommand

	ServiceManager service.Manager
	GetClient      GetClient
}

var _ flags.Commander = Login{}

func (c Login) Execute(args []string) error {
	rawEnvService, err := c.ServiceManager.Get("env")
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

	svc, err := c.ServiceManager.Get(authServiceType)
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

	env, err := c.Runtime.GetEnvironment()
	if err != nil {
		return bosherr.WrapError(err, "Getting environment state")
	}

	env.Auth = &config.EnvironmentAuthState{
		Type:    authServiceType,
		Options: auth,
	}

	configManager, err := c.Runtime.GetConfigManager()
	if err != nil {
		return bosherr.WrapError(err, "Getting config manager")
	}

	err = configManager.SetEnvironment(env)
	if err != nil {
		return bosherr.WrapError(err, "Updating environment")
	}

	// show confirmation of the new user

	ui := c.Runtime.GetUI()

	client, err := c.GetClient()
	if err != nil {
		return bosherr.WrapError(err, "Getting client")
	}

	authInfo, err := client.GetInfo()
	if err != nil {
		return bosherr.WrapError(err, "Getting remote authentication info")
	}

	if authInfo.ID == "" {
		return errors.New("Failed to use authentication credentials")
	}

	ui.PrintBlock(fmt.Sprintf("Successfully logged in as %s\n", authInfo.ID))

	return nil
}
