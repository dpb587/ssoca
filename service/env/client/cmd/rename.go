package cmd

import (
	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"

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
		return errors.Wrap(err, "Getting environment")
	}

	env.Alias = c.Args.Name

	configManager, err := c.Runtime.GetConfigManager()
	if err != nil {
		return errors.Wrap(err, "Getting state manager")
	}

	err = configManager.SetEnvironment(env)
	if err != nil {
		return errors.Wrap(err, "Setting environment")
	}

	return nil
}
