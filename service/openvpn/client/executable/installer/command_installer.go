package installer

import (
	"io"
	"os/exec"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/sirupsen/logrus"
)

type CommandInstaller struct {
	CmdRunner boshsys.CmdRunner // TODO use if Cmd can stream

	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader

	Sudo bool
	Exec string
	Args []string
}

func (i *CommandInstaller) Install(_ logrus.FieldLogger) error {
	executable, args := i.Exec, i.Args
	if i.Sudo {
		_, err := exec.LookPath("sudo")
		if err == nil {
			args = append([]string{executable}, args...)
			executable = "sudo"
		}
	}

	cmd := exec.Command(executable, args...)
	cmd.Stdin = i.Stdin
	cmd.Stdout = i.Stdout
	cmd.Stderr = i.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
