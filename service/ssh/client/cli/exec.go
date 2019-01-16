package cli

import (
	"github.com/dpb587/ssoca/cli/errors"
	clientcmd "github.com/dpb587/ssoca/client/cmd"
	svc "github.com/dpb587/ssoca/service/ssh/client"
	"github.com/jessevdk/go-flags"
)

type Exec struct {
	*clientcmd.ServiceCommand `no-flag:"true"`
	clientcmd.InteractiveAuthCommand

	Exec      string   `long:"exec" description:"Path to the ssh binary"`
	ExtraOpts []string `long:"opt" description:"Additional option to pass to ssh"` // @todo
	Args      ExecArgs `positional-args:"true" optional:"true"`

	serviceFactory svc.ServiceFactory
}

var _ flags.Commander = Exec{}

type ExecArgs struct {
	Host string `positional-arg-name:"HOST"`
}

func (c Exec) Execute(_ []string) error {
	service := c.serviceFactory.New(c.ServiceName)

	exit, err := service.Execute(svc.ExecuteOptions{
		Exec:      c.Exec,
		ExtraArgs: c.ExtraOpts,
		Host:      c.Args.Host,
	})

	if err != nil {
		return err
	}

	return errors.Exit{Code: exit}
}
