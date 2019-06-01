package client

import (
	"io/ioutil"
	"os"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/pkg/errors"
)

type ExecuteOptions struct {
	SkipAuthRetry bool

	Quiet      bool
	RemoteFile string

	ExtraArgs []string
}

func (s Service) Execute(opts ExecuteOptions) error {
	wd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "Getting working directory")
	}

	tmpfile, err := ioutil.TempFile(wd, ".ssoca-exec-")
	if err != nil {
		return errors.Wrap(err, "Creating temp file")
	}

	defer os.RemoveAll(tmpfile.Name())

	err = tmpfile.Close()
	if err != nil {
		return errors.Wrap(err, "Closing tempfile")
	}

	err = s.Get(GetOptions{
		SkipAuthRetry: opts.SkipAuthRetry,
		Quiet:         opts.Quiet,
		RemoteFile:    opts.RemoteFile,
		LocalFile:     tmpfile.Name(),
	})
	if err != nil {
		return errors.Wrap(err, "Getting file")
	}

	err = os.Chmod(tmpfile.Name(), 0700)
	if err != nil {
		return errors.Wrap(err, "Chmoding tempfile")
	}

	configManager, err := s.runtime.GetConfigManager()
	if err != nil {
		return errors.Wrap(err, "Getting config manager")
	}

	cmd := boshsys.Command{
		Name: tmpfile.Name(),
		Args: opts.ExtraArgs,
		Env: map[string]string{
			"SSOCA_CONFIG":      configManager.GetSource(),
			"SSOCA_ENVIRONMENT": s.runtime.GetEnvironmentName(),
		},
		Stdin:  s.runtime.GetStdin(),
		Stdout: s.runtime.GetStdout(),
		Stderr: s.runtime.GetStderr(),
	}

	_, _, _, err = s.cmdRunner.RunComplexCommand(cmd)
	if err != nil {
		return errors.Wrap(err, "Executing")
	}

	return nil
}
