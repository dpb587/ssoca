package cmd

import (
	boshtbl "github.com/cloudfoundry/bosh-cli/ui/table"
	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"

	clientcmd "github.com/dpb587/ssoca/client/cmd"
	svcapi "github.com/dpb587/ssoca/service/env/api"
)

type Services struct {
	*clientcmd.ServiceCommand `no-flag:"true"`

	GetClient GetClient
}

var _ flags.Commander = Info{}

func (c Services) Execute(_ []string) error {
	client, err := c.GetClient()
	if err != nil {
		return errors.Wrap(err, "getting client")
	}

	info, err := client.GetInfo()
	if err != nil {
		return errors.Wrap(err, "getting remote environment info")
	}

	table := boshtbl.Table{
		Header: []boshtbl.Header{
			{Title: "Name"},
			{Title: "Type"},
			{Title: "Metadata"},
		},
		Rows: [][]boshtbl.Value{},
	}

	// for old server + new client, avoid showing duplicate auth services
	deprecatedShowAuth := true

	for _, service := range info.Services {
		if service.Name == "auth" {
			deprecatedShowAuth = false
		}

		c.appendServiceRow(&table, service)
	}

	if deprecatedShowAuth && info.Auth != nil {
		info.Auth.Name = "auth"
		c.appendServiceRow(&table, *info.Auth)
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
