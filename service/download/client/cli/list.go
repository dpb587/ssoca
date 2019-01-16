package cli

import (
	boshtbl "github.com/cloudfoundry/bosh-cli/ui/table"
	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"

	clientcmd "github.com/dpb587/ssoca/client/cmd"
	svc "github.com/dpb587/ssoca/service/download/client"
)

type List struct {
	*clientcmd.ServiceCommand `no-flag:"true"`
	clientcmd.InteractiveAuthCommand

	serviceFactory svc.ServiceFactory
}

var _ flags.Commander = List{}

func (c List) Execute(_ []string) error {
	service := c.serviceFactory.New(c.ServiceName)

	files, err := service.List(svc.ListOptions{
		SkipAuthRetry: c.SkipAuthRetry,
	})
	if err != nil {
		return errors.Wrap(err, "Listing files")
	}

	table := boshtbl.Table{
		Header: []boshtbl.Header{
			{Title: "File"},
		},
	}

	for _, file := range files {
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
