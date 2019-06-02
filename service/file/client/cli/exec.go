package cli

import (
	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"

	clientcmd "github.com/dpb587/ssoca/client/cmd"
	svc "github.com/dpb587/ssoca/service/file/client"
)

type Exec struct {
	*clientcmd.ServiceCommand `no-flag:"true"`
	clientcmd.InteractiveAuthCommand

	serviceFactory svc.ServiceFactory

	Args ExecArgs `positional-args:"true"`
}

var _ flags.Commander = Exec{}

type ExecArgs struct {
	File  string   `positional-arg-name:"FILE" description:"File name" required:"true"`
	Extra []string `positional-arg-name:"EXTRA" description:"Additional arguments to pass"`
}

func (c Exec) Execute(_ []string) error {
	service := c.serviceFactory.New(c.ServiceName)

	err := service.Execute(svc.ExecuteOptions{
		SkipAuthRetry: c.SkipAuthRetry,
		RemoteFile:    c.Args.File,
		ExtraArgs:     c.Args.Extra,
	})
	if err != nil {
		return errors.Wrap(err, "executing file")
	}

	return nil
}
