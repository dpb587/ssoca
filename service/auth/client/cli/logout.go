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
		return errors.Wrap(err, "getting environment state")
	}

	if env.Auth == nil {
		return nil
	}

	authServiceType := env.Auth.Type

	svc, err := c.ServiceManager.Get(authServiceType)
	if err != nil {
		return errors.Wrap(err, "loading auth service")
	}

	authService, ok := svc.(service.AuthService)
	if !ok {
		return fmt.Errorf("cannot authenticate with service: %s", authServiceType)
	}

	err = authService.AuthLogout()
	if err != nil {
		return errors.Wrap(err, "unauthenticating")
	}

	env.Auth = nil

	configManager, err := c.Runtime.GetConfigManager()
	if err != nil {
		return errors.Wrap(err, "getting state manager")
	}

	err = configManager.SetEnvironment(env)
	if err != nil {
		return errors.Wrap(err, "updating environment state")
	}

	ui := c.Runtime.GetUI()

	ui.PrintBlock([]byte("Successfully logged out\n"))

	return nil
}
