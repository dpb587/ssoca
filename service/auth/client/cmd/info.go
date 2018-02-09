package cmd

import (
	boshtbl "github.com/cloudfoundry/bosh-cli/ui/table"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	clientcmd "github.com/dpb587/ssoca/client/cmd"
	"github.com/jessevdk/go-flags"
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
		return bosherr.WrapError(err, "Getting remote authentication info")
	}

	table := boshtbl.Table{}

	if info.ID != "" {
		table.Rows = append(
			table.Rows,
			[]boshtbl.Value{
				boshtbl.NewValueString("Authenticated"),
				boshtbl.NewValueBool(true),
			},
		)

		table.Rows = append(
			table.Rows,
			[]boshtbl.Value{
				boshtbl.NewValueString("ID"),
				boshtbl.NewValueString(info.ID),
			},
		)

		for _, k := range info.Groups {
			table.Rows = append(
				table.Rows,
				[]boshtbl.Value{
					boshtbl.NewValueString("Group"),
					boshtbl.NewValueString(k),
					boshtbl.NewValueInterface(nil),
				},
			)
		}

		for k, v := range info.Attributes {
			table.Rows = append(
				table.Rows,
				[]boshtbl.Value{
					boshtbl.NewValueString("Attribute"),
					boshtbl.NewValueString(k),
					boshtbl.NewValueString(v),
				},
			)
		}
	} else {
		table.Rows = append(
			table.Rows,
			[]boshtbl.Value{
				boshtbl.NewValueString("Authenticated"),
				boshtbl.NewValueBool(false),
			},
		)
	}

	c.Runtime.GetUI().PrintTable(table)

	return nil
}
