package cmd

import (
	boshtbl "github.com/cloudfoundry/bosh-cli/ui/table"
	clientcmd "github.com/dpb587/ssoca/client/cmd"
	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
)

type List struct {
	clientcmd.ServiceCommand
}

var _ flags.Commander = List{}

func (c List) Execute(_ []string) error {
	table := boshtbl.Table{
		Header: []boshtbl.Header{
			{Title: "URL"},
			{Title: "Alias"},
		},
	}

	configManager, err := c.Runtime.GetConfigManager()
	if err != nil {
		return errors.Wrap(err, "Getting state manager")
	}

	envs, err := configManager.GetEnvironments()
	if err != nil {
		return errors.Wrap(err, "Getting environments")
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
