// +build !windows

package client

import boshsys "github.com/cloudfoundry/bosh-utils/system"

func executeRewriteCommand(cmd boshsys.Command) boshsys.Command {
	return cmd
}
