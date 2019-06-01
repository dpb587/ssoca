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
			"Chocolatey": &CommandInstaller{
				Exec:   "choco",
				Args:   []string{"install", "openvpn"},
				Stdin:  runtime.GetStdin(),
				Stdout: runtime.GetStdout(),
				Stderr: runtime.GetStderr(),
			},
			"fallback steps": &MessageInstaller{
				Output: runtime.GetStdout(),
				Message: `Try downloading and installing the official OpenVPN GUI...

    https://openvpn.net/index.php/open-source/downloads.html)`,
			},
		},
	}
}
