package cmd

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/jessevdk/go-flags"

	clientcmd "github.com/dpb587/ssoca/client/cmd"

	"github.com/dpb587/ssoca/client/config"
)

type Add struct {
	clientcmd.ServiceCommand

	Args              AddArgs `positional-args:"true"`
	CACertificatePath string  `long:"ca-cert" description:"Environment CA certificate path"`

	FS boshsys.FileSystem
}

var _ flags.Commander = Add{}

type AddArgs struct {
	URL string `positional-arg-name:"URI" description:"Environment URL"`
}

func (c Add) Execute(_ []string) error {
	env := config.EnvironmentState{
		Alias: c.Runtime.GetEnvironmentName(),
		URL:   c.Args.URL,
	}

	if c.CACertificatePath != "" {
		absPath, err := c.FS.ExpandPath(c.CACertificatePath)
		if err != nil {
			return bosherr.WrapError(err, "Expanding path")
		}

		cacert, err := c.FS.ReadFileString(absPath)
		if err != nil {
			return bosherr.WrapError(err, "Reading file")
		}

		env.CACertificate = cacert
	}

	configManager, err := c.Runtime.GetConfigManager()
	if err != nil {
		return bosherr.WrapError(err, "Getting state manager")
	}

	err = configManager.SetEnvironment(env)
	if err != nil {
		return bosherr.WrapError(err, "Adding environment")
	}

	return nil
}
