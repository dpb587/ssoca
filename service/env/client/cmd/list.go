package cmd

import (
	boshtbl "github.com/cloudfoundry/bosh-cli/ui/table"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	clientcmd "github.com/dpb587/ssoca/client/cmd"
	"github.com/jessevdk/go-flags"
)

type List struct {
	clientcmd.ServiceCommand
}

var _ flags.Commander = List{}

func (c List) Execute(args []string) error {
	table := boshtbl.Table{
		Header: []string{
			"URL",
			"Alias",
		},
	}

	configManager, err := c.Runtime.GetConfigManager()
	if err != nil {
		return bosherr.WrapError(err, "Getting state manager")
	}

	envs, err := configManager.GetEnvironments()
	if err != nil {
		return bosherr.WrapError(err, "Getting environments")
	}

	for _, env := range envs {
		table.Rows = append(
			table.Rows,
			[]boshtbl.Value{
				boshtbl.NewValueString(env.URL),
				boshtbl.NewValueString(env.Alias),
			},
		)
	}

	c.Runtime.GetUI().PrintTable(table)

	return nil
}
