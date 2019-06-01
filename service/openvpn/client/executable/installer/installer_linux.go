package installer

import (
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/dpb587/ssoca/client"
	"github.com/dpb587/ssoca/service/openvpn/client/executable/finder"
)

func New(runtime client.Runtime, cmdRunner boshsys.CmdRunner, f finder.Finder) Installer {
	return &MultiInstaller{
		Name:   "openvpn",
		Finder: f,
		Installers: map[string]Installer{
			"apt (sudo)": &CommandInstaller{
				CmdRunner: cmdRunner,
				Sudo:      true,
				Exec:      "apt",
				Args:      []string{"install", "-y", "openvpn"},
				Stdin:     runtime.GetStdin(),
				Stdout:    runtime.GetStdout(),
				Stderr:    runtime.GetStderr(),
			},
			"yum (sudo)": &CommandInstaller{
				CmdRunner: cmdRunner,
				Sudo:      true,
				Exec:      "yum",
				Args:      []string{"install", "-y", "openvpn"},
				Stdin:     runtime.GetStdin(),
				Stdout:    runtime.GetStdout(),
				Stderr:    runtime.GetStderr(),
			},
		},
	}
}
