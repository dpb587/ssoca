package cmd

import (
	"fmt"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/jessevdk/go-flags"

	clientcmd "github.com/dpb587/ssoca/client/cmd"
	"github.com/dpb587/ssoca/client/service"
)

type Logout struct {
	clientcmd.ServiceCommand
	ServiceManager service.Manager
}

var _ flags.Commander = Logout{}

func (c Logout) Execute(args []string) error {
	env, err := c.Runtime.GetEnvironment()
	if err != nil {
		return bosherr.WrapError(err, "Getting environment state")
	}

	if env.Auth == nil {
		return nil
	}

	authServiceType := env.Auth.Type

	svc, err := c.ServiceManager.Get(authServiceType)
	if err != nil {
		return bosherr.WrapError(err, "Loading auth service")
	}

	authService, ok := svc.(service.AuthService)
	if !ok {
		return fmt.Errorf("Cannot authenticate with service: %s", authServiceType)
	}

	err = authService.AuthLogout()
	if err != nil {
		return bosherr.WrapError(err, "Unauthenticating")
	}

	env.Auth = nil

	configManager, err := c.Runtime.GetConfigManager()
	if err != nil {
		return bosherr.WrapError(err, "Getting state manager")
	}

	err = configManager.SetEnvironment(env)
	if err != nil {
		return bosherr.WrapError(err, "Updating environment state")
	}

	ui := c.Runtime.GetUI()

	ui.PrintBlock("Successfully logged out\n")

	return nil
}
