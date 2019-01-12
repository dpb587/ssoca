package cmd

import (
	boshtbl "github.com/cloudfoundry/bosh-cli/ui/table"
	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"

	clientcmd "github.com/dpb587/ssoca/client/cmd"
)

type List struct {
	clientcmd.ServiceCommand
	clientcmd.InteractiveAuthCommand

	GetClient GetClient
}

var _ flags.Commander = List{}

func (c List) Execute(_ []string) error {
	client, err := c.GetClient(c.ServiceName, c.SkipAuthRetry)
	if err != nil {
		return errors.Wrap(err, "Getting client")
	}

	list, err := client.GetList()
	if err != nil {
		return errors.Wrap(err, "Getting remote environment info")
	}

	table := boshtbl.Table{
		Header: []boshtbl.Header{
			{Title: "File"},
		},
	}

	for _, file := range list.Files {
		table.Rows = append(
			table.Rows,
			[]boshtbl.Value{
				boshtbl.NewValueString(file.Name),
			},
		)
	}

	c.Runtime.GetUI().PrintTable(table)

	return nil
}
