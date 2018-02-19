package cmd

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/jessevdk/go-flags"

	clientcmd "github.com/dpb587/ssoca/client/cmd"
)

type Rename struct {
	clientcmd.ServiceCommand

	GetClient GetClient

	Args RenameArgs `positional-args:"true"`
}

var _ flags.Commander = Rename{}

type RenameArgs struct {
	Name string `positional-arg-name:"NEW-NAME" description:"New environment name"`
}

func (c Rename) Execute(_ []string) error {
	env, err := c.Runtime.GetEnvironment()
	if err != nil {
		return bosherr.WrapError(err, "Getting environment")
	}

	env.Alias = c.Args.Name

	configManager, err := c.Runtime.GetConfigManager()
	if err != nil {
		return bosherr.WrapError(err, "Getting state manager")
	}

	err = configManager.SetEnvironment(env)
	if err != nil {
		return bosherr.WrapError(err, "Setting environment")
	}

	return nil
}
