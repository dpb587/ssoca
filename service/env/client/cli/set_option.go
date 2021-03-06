package cli

import (
	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"

	clientcmd "github.com/dpb587/ssoca/client/cmd"
	yaml "gopkg.in/yaml.v2"
)

type SetOption struct {
	*clientcmd.ServiceCommand `no-flag:"true"`

	Args              SetOptionArgs `positional-args:"true"`
	CACertificatePath string        `long:"ca-cert" description:"Environment CA certificate path"`
}

var _ flags.Commander = SetOption{}

type SetOptionArgs struct {
	Name  string `positional-arg-name:"NAME" description:"Client option name"`
	Value string `positional-arg-name:"VALUE" description:"Client option value (parsed as YAML)"`
}

func (c SetOption) Execute(_ []string) error {
	env, err := c.Runtime.GetEnvironment()
	if err != nil {
		return errors.Wrap(err, "getting environment")
	}

	rawValue := c.Args.Value
	var parsedValue interface{}

	err = yaml.Unmarshal([]byte(rawValue), &parsedValue)
	if err != nil {
		return errors.Wrap(err, "unmarshaling YAML value")
	}

	env.SetOption(c.Args.Name, parsedValue)

	configManager, err := c.Runtime.GetConfigManager()
	if err != nil {
		return errors.Wrap(err, "getting state manager")
	}

	err = configManager.SetEnvironment(env)
	if err != nil {
		return errors.Wrap(err, "setting environment")
	}

	return nil
}
