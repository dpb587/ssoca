package cli

import (
	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"

	clientcmd "github.com/dpb587/ssoca/client/cmd"
	svc "github.com/dpb587/ssoca/service/file/client"
)

type Get struct {
	*clientcmd.ServiceCommand `no-flag:"true"`
	clientcmd.InteractiveAuthCommand

	serviceFactory svc.ServiceFactory

	Args GetArgs `positional-args:"true"`
}

var _ flags.Commander = Get{}

type GetArgs struct {
	File       string `positional-arg-name:"FILE" description:"File name" required:"true"`
	TargetFile string `positional-arg-name:"TARGET-FILE" description:"Target path to write download (use '-' for STDOUT)"`
}

func (c Get) Execute(_ []string) error {
	service := c.serviceFactory.New(c.ServiceName)

	err := service.Get(svc.GetOptions{
		SkipAuthRetry: c.SkipAuthRetry,
		RemoteFile:    c.Args.File,
		LocalFile:     c.Args.TargetFile,
	})
	if err != nil {
		return errors.Wrap(err, "Getting file")
	}

	return nil
}
