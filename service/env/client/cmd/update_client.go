package cmd

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/cheggaaa/pb"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/inconshreveable/go-update"
	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"

	clientcmd "github.com/dpb587/ssoca/client/cmd"
	downloadhttpclient "github.com/dpb587/ssoca/service/file/httpclient"
)

type UpdateClient struct {
	*clientcmd.ServiceCommand `no-flag:"true"`
	clientcmd.InteractiveAuthCommand

	FS        boshsys.FileSystem
	SsocaExec string
	CmdRunner boshsys.CmdRunner

	GetClient         GetClient
	GetDownloadClient GetDownloadClient
}

var _ flags.Commander = UpdateClient{}

func (c UpdateClient) Execute(_ []string) error {
	client, err := c.GetClient()
	if err != nil {
		return errors.Wrap(err, "getting client")
	}

	info, err := client.GetInfo()
	if err != nil {
		return errors.Wrap(err, "getting remote environment info")
	}

	if info.Env.UpdateService == "" {
		return errors.New("environment does not provide a client update service")
	}

	downloadClient, err := c.GetDownloadClient(info.Env.UpdateService, c.SkipAuthRetry)
	if err != nil {
		return errors.Wrap(err, "getting download client")
	}

	metadata, err := downloadClient.GetMetadata()
	if err != nil {
		return errors.Wrap(err, "getting download metadata")
	}

	version, ok := metadata.Metadata["version"]
	if !ok {
		return errors.New("environment does not advertise the client version")
	}

	if version == c.Runtime.GetVersion().Semver {
		return nil
	}

	files, err := downloadClient.GetList()
	if err != nil {
		return errors.Wrap(err, "listing client files")
	}

	var found string

	for _, file := range files.Files {
		if !strings.Contains(file.Name, "ssoca-client-") {
			continue
		} else if !strings.Contains(file.Name, fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)) {
			continue
		} else if found != "" {
			return fmt.Errorf("multiple clients were found: %s, %s", found, file.Name)
		}

		found = file.Name
	}

	if found == "" {
		return fmt.Errorf("unable to find client (%s, %s)", runtime.GOOS, runtime.GOARCH)
	}

	executable := c.SsocaExec
	if executable == "" {
		executable = "ssoca"
	}

	executable, err = exec.LookPath(executable)
	if err != nil {
		return errors.Wrap(err, "expanding path")
	}

	err = c.update(downloadClient, executable, found)
	if err != nil {
		return errors.Wrap(err, "updating binary")
	}

	_, _, exit, err := c.CmdRunner.RunComplexCommand(boshsys.Command{
		Name:   executable,
		Args:   []string{"version"},
		Stderr: c.Runtime.GetStderr(),
		Stdout: c.Runtime.GetStdout(),
	})
	if err != nil {
		return errors.Wrap(err, "verifying updated binary")
	} else if exit != 0 {
		return fmt.Errorf("unexpected exit from updated binary: %d", exit)
	}

	return nil
}

func (c UpdateClient) update(downloadClient downloadhttpclient.Client, executable string, fileName string) error {
	tmpfile, err := c.FS.TempFile("ssoca-update-client-")
	if err != nil {
		return errors.Wrap(err, "creating temporary file for download")
	}

	defer tmpfile.Close()

	downloadStatus := pb.New(0).SetRefreshRate(250 * time.Millisecond).SetWidth(80)
	downloadStatus.Output = c.Runtime.GetStderr()
	downloadStatus.ShowPercent = false

	err = downloadClient.Download(fileName, tmpfile, downloadStatus)
	if err != nil {
		return errors.Wrap(err, "downloading file")
	}

	_, err = tmpfile.Seek(0, 0)
	if err != nil {
		return errors.Wrap(err, "rewinding download")
	}

	err = update.Apply(tmpfile, update.Options{TargetPath: executable})
	if err != nil {
		return errors.Wrap(err, "updating file")
	}

	return nil
}
