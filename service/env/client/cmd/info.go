package cmd

import (
	boshtbl "github.com/cloudfoundry/bosh-cli/ui/table"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	clientcmd "github.com/dpb587/ssoca/client/cmd"
	svcapi "github.com/dpb587/ssoca/service/env/api"
)

type Info struct {
	clientcmd.ServiceCommand

	GetClient GetClient
}

func (c Info) Execute(args []string) error {
	client, err := c.GetClient()
	if err != nil {
		return bosherr.WrapError(err, "Getting client")
	}

	info, err := client.GetInfo()
	if err != nil {
		return bosherr.WrapError(err, "Getting remote environment info")
	}

	table := boshtbl.Table{
		HeaderVals: []boshtbl.Value{
			boshtbl.NewValueString("Service"),
			boshtbl.NewValueString("Type"),
			boshtbl.NewValueString("Metadata"),
		},
		Rows: [][]boshtbl.Value{},
	}

	info.Auth.Name = "auth"
	c.appendServiceRow(&table, info.Auth)

	for _, service := range info.Services {
		c.appendServiceRow(&table, service)
	}

	c.Runtime.GetUI().PrintTable(table)

	return nil
}

func (c Info) appendServiceRow(table *boshtbl.Table, service svcapi.InfoServiceResponse) {
	table.Rows = append(
		table.Rows,
		[]boshtbl.Value{
			boshtbl.NewValueString(service.Name),
			boshtbl.NewValueString(service.Type),
			boshtbl.NewValueStrings([]string{}),
		},
	)
}
