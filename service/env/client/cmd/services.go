package cmd

import (
	boshtbl "github.com/cloudfoundry/bosh-cli/ui/table"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/jessevdk/go-flags"

	clientcmd "github.com/dpb587/ssoca/client/cmd"
	svcapi "github.com/dpb587/ssoca/service/env/api"
)

type Services struct {
	clientcmd.ServiceCommand

	GetClient GetClient
}

var _ flags.Commander = Info{}

func (c Services) Execute(_ []string) error {
	client, err := c.GetClient()
	if err != nil {
		return bosherr.WrapError(err, "Getting client")
	}

	info, err := client.GetInfo()
	if err != nil {
		return bosherr.WrapError(err, "Getting remote environment info")
	}

	table := boshtbl.Table{
		Header: []boshtbl.Header{
			{Title: "Name"},
			{Title: "Type"},
			{Title: "Metadata"},
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

func (c Services) appendServiceRow(table *boshtbl.Table, service svcapi.InfoServiceResponse) {
	table.Rows = append(
		table.Rows,
		[]boshtbl.Value{
			boshtbl.NewValueString(service.Name),
			boshtbl.NewValueString(service.Type),
			boshtbl.NewValueStrings([]string{}),
		},
	)
}
