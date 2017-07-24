package cmd

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/jessevdk/go-flags"

	clientcmd "github.com/dpb587/ssoca/client/cmd"
	yaml "gopkg.in/yaml.v2"
)

type SetOption struct {
	clientcmd.ServiceCommand

	Args              SetOptionArgs `positional-args:"true"`
	CACertificatePath string        `long:"ca-cert" description:"Environment CA certificate path"`

	FS boshsys.FileSystem
}

var _ flags.Commander = SetOption{}

type SetOptionArgs struct {
	Name  string `positional-arg-name:"NAME" description:"Client option name"`
	Value string `positional-arg-name:"VALUE" description:"Client option value (parsed as YAML)"`
}

func (c SetOption) Execute(args []string) error {
	env, err := c.Runtime.GetEnvironment()
	if err != nil {
		return bosherr.WrapError(err, "Getting environment")
	}

	rawValue := c.Args.Value
	var parsedValue interface{}

	err = yaml.Unmarshal([]byte(rawValue), &parsedValue)
	if err != nil {
		return bosherr.WrapError(err, "Unmarshaling YAML value")
	}

	env.SetOption(c.Args.Name, parsedValue)

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
