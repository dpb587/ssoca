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
			"Homebrew": &CommandInstaller{
				CmdRunner: cmdRunner,
				Exec:      "brew",
				Args:      []string{"install", "openvpn"},
				Stdin:     runtime.GetStdin(),
				Stdout:    runtime.GetStdout(),
				Stderr:    runtime.GetStderr(),
			},
			"fallback steps": &MessageInstaller{
				Output: runtime.GetStdout(),
				Message: `
Try installing via Homebrew...

    brew install openvpn

Alternatively, the following applications will also install openvpn...

 * Tunnelblick (https://tunnelblick.net/)
 * Shimo (https://www.shimovpn.com/)
 * Viscosity (https://www.sparklabs.com/viscosity/)
`,
			},
		},
	}
}
