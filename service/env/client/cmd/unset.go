package cmd

import (
	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"

	clientcmd "github.com/dpb587/ssoca/client/cmd"
)

type Unset struct {
	*clientcmd.ServiceCommand `no-flag:"true"`
}

var _ flags.Commander = Unset{}

func (c Unset) Execute(_ []string) error {
	configManager, err := c.Runtime.GetConfigManager()
	if err != nil {
		return errors.Wrap(err, "Getting state manager")
	}

	err = configManager.UnsetEnvironment(c.Runtime.GetEnvironmentName())
	if err != nil {
		return errors.Wrap(err, "Unsetting environment")
	}

	return nil
}
