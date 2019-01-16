package cli

import (
	"fmt"

	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"

	clientcmd "github.com/dpb587/ssoca/client/cmd"
	"github.com/dpb587/ssoca/client/service"
)

type Logout struct {
	*clientcmd.ServiceCommand `no-flag:"true"`

	ServiceManager service.Manager
}

var _ flags.Commander = Logout{}

func (c Logout) Execute(_ []string) error {
	env, err := c.Runtime.GetEnvironment()
	if err != nil {
		return errors.Wrap(err, "Getting environment state")
	}

	if env.Auth == nil {
		return nil
	}

	authServiceType := env.Auth.Type

	svc, err := c.ServiceManager.Get(authServiceType)
	if err != nil {
		return errors.Wrap(err, "Loading auth service")
	}

	authService, ok := svc.(service.AuthService)
	if !ok {
		return fmt.Errorf("Cannot authenticate with service: %s", authServiceType)
	}

	err = authService.AuthLogout()
	if err != nil {
		return errors.Wrap(err, "Unauthenticating")
	}

	env.Auth = nil

	configManager, err := c.Runtime.GetConfigManager()
	if err != nil {
		return errors.Wrap(err, "Getting state manager")
	}

	err = configManager.SetEnvironment(env)
	if err != nil {
		return errors.Wrap(err, "Updating environment state")
	}

	ui := c.Runtime.GetUI()

	ui.PrintBlock([]byte("Successfully logged out\n"))

	return nil
}
