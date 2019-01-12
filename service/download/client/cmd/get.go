package cmd

import (
	"os"
	"time"

	"github.com/cheggaaa/pb"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	clientcmd "github.com/dpb587/ssoca/client/cmd"
	"github.com/jessevdk/go-flags"
)

type Get struct {
	clientcmd.ServiceCommand
	clientcmd.InteractiveAuthCommand

	GetClient GetClient
	FS        boshsys.FileSystem

	Args GetArgs `positional-args:"true"`
}

var _ flags.Commander = Get{}

type GetArgs struct {
	File       string `positional-arg-name:"FILE" description:"File name" required:"true"`
	TargetFile string `positional-arg-name:"TARGET-FILE" description:"Target path to write download"`
}

func (c Get) Execute(_ []string) error {
	client, err := c.GetClient(c.ServiceName, c.SkipAuthRetry)
	if err != nil {
		return bosherr.WrapError(err, "Getting client")
	}

	filePath := c.Args.TargetFile

	if filePath == "" {
		filePath = c.Args.File
	}

	file, err := c.FS.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0600)
	if err != nil {
		return bosherr.WrapError(err, "Opening file")
	}

	downloadStatus := pb.New(0).SetRefreshRate(250 * time.Millisecond).SetWidth(80)
	downloadStatus.Output = c.Runtime.GetStderr()
	downloadStatus.ShowPercent = false

	return client.Download(c.Args.File, file, downloadStatus)
}
