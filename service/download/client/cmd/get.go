package cmd

import (
	"io"
	"os"
	"time"

	"github.com/cheggaaa/pb"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"

	clientcmd "github.com/dpb587/ssoca/client/cmd"
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
	TargetFile string `positional-arg-name:"TARGET-FILE" description:"Target path to write download (use '-' for STDOUT)"`
}

func (c Get) Execute(_ []string) error {
	client, err := c.GetClient(c.ServiceName, c.SkipAuthRetry)
	if err != nil {
		return errors.Wrap(err, "Getting client")
	}

	filePath := c.Args.TargetFile

	if filePath == "" {
		filePath = c.Args.File
	}

	var file io.ReadWriteSeeker

	if filePath == "-" {
		fileTemp, err := c.FS.TempFile("ssoca-download-stdout")
		if err != nil {
			return errors.Wrap(err, "Creating temp file for stdout")
		}

		defer os.RemoveAll(fileTemp.Name())

		file = fileTemp
	} else {
		file, err = c.FS.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0600)
		if err != nil {
			return errors.Wrap(err, "Opening file")
		}
	}

	downloadStatus := pb.New(0).SetRefreshRate(250 * time.Millisecond).SetWidth(80)
	downloadStatus.Output = c.Runtime.GetStderr()
	downloadStatus.ShowPercent = false

	err = client.Download(c.Args.File, file, downloadStatus)
	if err != nil {
		return err
	}

	if filePath != "-" {
		return nil
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return errors.Wrap(err, "Seeking downloaded temp file")
	}

	// TODO stdout should use runtime.GetStdout()?
	_, err = io.Copy(os.Stdout, file)
	if err != nil {
		return errors.Wrap(err, "Writing download to STDOUT")
	}

	return nil
}
