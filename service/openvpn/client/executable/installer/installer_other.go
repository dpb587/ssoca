// +build !darwin
// +build !linux
// +build !windows

package installer

import (
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/dpb587/ssoca/client"
	"github.com/dpb587/ssoca/service/openvpn/client/executable/finder"
)

func New(runtime client.Runtime, cmdRunner boshsys.CmdRunner, f finder.Finder) Installer {
	return &MultiInstaller{}
}
