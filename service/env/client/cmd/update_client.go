package cmd

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/cheggaaa/pb"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/inconshreveable/go-update"
	"github.com/jessevdk/go-flags"

	clientcmd "github.com/dpb587/ssoca/client/cmd"
	downloadhttpclient "github.com/dpb587/ssoca/service/download/httpclient"
)

type UpdateClient struct {
	clientcmd.ServiceCommand
	clientcmd.InteractiveAuthCommand

	FS        boshsys.FileSystem
	SsocaExec string

	GetClient         GetClient
	GetDownloadClient GetDownloadClient
}

var _ flags.Commander = UpdateClient{}

func (c UpdateClient) Execute(_ []string) error {
	client, err := c.GetClient()
	if err != nil {
		return bosherr.WrapError(err, "Getting client")
	}

	info, err := client.GetInfo()
	if err != nil {
		return bosherr.WrapError(err, "Getting remote environment info")
	}

	if info.Env.UpdateService == "" {
		return errors.New("Environment does not provide a client update service")
	}

	downloadClient, err := c.GetDownloadClient(info.Env.UpdateService, c.SkipAuthRetry)
	if err != nil {
		return bosherr.WrapError(err, "Getting download client")
	}

	metadata, err := downloadClient.GetMetadata()
	if err != nil {
		return bosherr.WrapError(err, "Getting download metadata")
	}

	version, ok := metadata.Metadata["version"]
	if !ok {
		return errors.New("Environment does not advertise the client version")
	}

	if version == c.Runtime.GetVersion().Semver {
		return nil
	}

	files, err := downloadClient.GetList()
	if err != nil {
		return bosherr.WrapError(err, "Listing client files")
	}

	var found string

	for _, file := range files.Files {
		if !strings.Contains(file.Name, "ssoca-client-") {
			continue
		} else if !strings.Contains(file.Name, fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)) {
			continue
		} else if found != "" {
			return fmt.Errorf("Multiple clients were found: %s, %s", found, file.Name)
		}

		found = file.Name
	}

	if found == "" {
		return fmt.Errorf("Unable to find client (%s, %s)", runtime.GOOS, runtime.GOARCH)
	}

	return c.update(downloadClient, found)
}

func (c UpdateClient) update(downloadClient downloadhttpclient.Client, fileName string) error {
	executable := c.SsocaExec
	if executable == "" {
		executable = "ssoca"
	}

	executable, err := exec.LookPath(executable)
	if err != nil {
		return bosherr.WrapError(err, "Expanding path")
	}

	tmpfile, err := c.FS.TempFile("ssoca-update-client-")
	if err != nil {
		return bosherr.WrapError(err, "Creating temporary file for download")
	}

	defer tmpfile.Close()

	downloadStatus := pb.New(0).SetRefreshRate(250 * time.Millisecond).SetWidth(80)
	downloadStatus.Output = c.Runtime.GetStderr()
	downloadStatus.ShowPercent = false

	err = downloadClient.Download(fileName, tmpfile, downloadStatus)
	if err != nil {
		return bosherr.WrapError(err, "Downloading file")
	}

	_, err = tmpfile.Seek(0, 0)
	if err != nil {
		return bosherr.WrapError(err, "Rewinding download")
	}

	err = update.Apply(tmpfile, update.Options{TargetPath: executable})
	if err != nil {
		return bosherr.WrapError(err, "Updating file")
	}

	return nil
}
