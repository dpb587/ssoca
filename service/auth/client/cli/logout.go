package cli

import (
	"context"
	"fmt"

	clientcmd "github.com/dpb587/ssoca/client/cmd"
	"github.com/dpb587/ssoca/client/service"
	globalservice "github.com/dpb587/ssoca/service"
	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
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

	authServiceType := globalservice.Type(env.Auth.Type)

	svc, err := c.ServiceManager.Get(authServiceType, env.Auth.Name)
	if err != nil {
		return errors.Wrap(err, "loading auth service")
	}

	authService, ok := svc.(service.AuthService)
	if !ok {
		return fmt.Errorf("cannot authenticate with service: %s", authServiceType)
	}

	err = authService.AuthLogout(context.Background())
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
