package cmd

import (
	boshtbl "github.com/cloudfoundry/bosh-cli/ui/table"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/jessevdk/go-flags"

	clientcmd "github.com/dpb587/ssoca/client/cmd"
)

type Info struct {
	clientcmd.ServiceCommand

	GetClient GetClient
}

var _ flags.Commander = Info{}

func (c Info) Execute(_ []string) error {
	client, err := c.GetClient()
	if err != nil {
		return bosherr.WrapError(err, "Getting client")
	}

	info, err := client.GetInfo()
	if err != nil {
		return bosherr.WrapError(err, "Getting remote environment info")
	}

	table := boshtbl.Table{}

	table.Rows = append(
		table.Rows,
		[]boshtbl.Value{
			boshtbl.NewValueString("Title"),
			boshtbl.NewValueString(info.Env.Title),
		},
	)

	table.Rows = append(
		table.Rows,
		[]boshtbl.Value{
			boshtbl.NewValueString("URL"),
			boshtbl.NewValueString(info.Env.URL),
		},
	)

	if info.Env.Banner != "" {
		table.Rows = append(
			table.Rows,
			[]boshtbl.Value{
				boshtbl.NewValueString("Banner"),
				boshtbl.NewValueString(info.Env.Banner),
			},
		)
	}

	c.Runtime.GetUI().PrintTable(table)

	return nil
}
