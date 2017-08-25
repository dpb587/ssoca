package client

import (
	"fmt"
	"strings"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

func executeRewriteCommand(cmd boshsys.Command) boshsys.Command {
	cmd.Args = append([]string{cmd.Name}, cmd.Args...)
	for i, arg := range cmd.Args {
		cmd.Args[i] = fmt.Sprintf("'%s'", arg)
	}

	cmd.Name = "-Command"
	cmd.Args = []string{"& " + strings.Join(cmd.Args, " ")}

	return cmd
}
