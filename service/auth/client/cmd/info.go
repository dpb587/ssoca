package cmd

import (
	"encoding/json"

	boshtbl "github.com/cloudfoundry/bosh-cli/ui/table"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	clientcmd "github.com/dpb587/ssoca/client/cmd"
)

type Info struct {
	clientcmd.ServiceCommand

	GetClient GetClient
}

func (c *Info) Execute(args []string) error {
	client, err := c.GetClient()
	if err != nil {
		return bosherr.WrapError(err, "Getting client")
	}

	info, err := client.GetInfo()
	if err != nil {
		return bosherr.WrapError(err, "Getting remote authentication info")
	}

	table := boshtbl.Table{
	// Rows: [][]boshtbl.Value{},
	}

	if info.Username != "" {
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
				boshtbl.NewValueString("Username"),
				boshtbl.NewValueString(info.Username),
			},
		)

		for k, v := range info.Attributes {
			marshal, err := json.Marshal(v)
			if err != nil {
				return bosherr.WrapErrorf(err, "Marshalling attribute %s", k)
			}

			table.Rows = append(
				table.Rows,
				[]boshtbl.Value{
					boshtbl.NewValueString("Attribute"),
					boshtbl.NewValueString(k),
					boshtbl.NewValueString(string(marshal)),
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
