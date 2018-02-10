package cmd

import (
	boshtbl "github.com/cloudfoundry/bosh-cli/ui/table"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	clierrors "github.com/dpb587/ssoca/cli/errors"
	clientcmd "github.com/dpb587/ssoca/client/cmd"
	"github.com/jessevdk/go-flags"
)

type Info struct {
	clientcmd.ServiceCommand

	GetClient GetClient

	Authenticated bool `long:"authenticated" description:"Show only whether the user is authenticated"`
	ID            bool `long:"id" description:"Show only the ID of the authenticated user"`
	Groups        bool `long:"groups" description:"Show only the groups of the authenticated user"`
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

	ui := c.Runtime.GetUI()

	if c.Authenticated {
		if info.ID != "" {
			ui.PrintBlock([]byte("true\n"))
		} else {
			ui.PrintBlock([]byte("false\n"))
		}
	} else if c.ID {
		ui.PrintBlock(append([]byte(info.ID), '\n'))
	} else if c.Groups {
		for _, k := range info.Groups {
			ui.PrintBlock(append([]byte(k), '\n'))
		}
	} else {
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

		ui.PrintTable(table)
	}

	if info.ID == "" {
		return clierrors.Exit{Code: 1}
	}

	return nil
}
