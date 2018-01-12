package cmd

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/jessevdk/go-flags"

	clientcmd "github.com/dpb587/ssoca/client/cmd"
)

type Unset struct {
	clientcmd.ServiceCommand
}

var _ flags.Commander = Unset{}

func (c Unset) Execute(_ []string) error {
	configManager, err := c.Runtime.GetConfigManager()
	if err != nil {
		return bosherr.WrapError(err, "Getting state manager")
	}

	err = configManager.UnsetEnvironment(c.Runtime.GetEnvironmentName())
	if err != nil {
		return bosherr.WrapError(err, "Unsetting environment")
	}

	return nil
}
