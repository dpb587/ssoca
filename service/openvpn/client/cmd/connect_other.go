// +build !windows

package cmd

import boshsys "github.com/cloudfoundry/bosh-utils/system"

func (c Connect) osCommand(cmd boshsys.Command) boshsys.Command {
	return cmd
}
