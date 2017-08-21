package cli

import (
	clientcmd "github.com/dpb587/ssoca/client/cmd"
	svc "github.com/dpb587/ssoca/service/openvpn/client"
	"github.com/jessevdk/go-flags"
)

type CreateTunnelblickProfile struct {
	clientcmd.ServiceCommand
	clientcmd.InteractiveAuthCommand

	SssocaExec string                       `long:"exec-ssoca" description:"Path to the ssoca binary"`
	Name       string                       `long:"name" description:"Specific file name to use for *.tblk"`
	Args       createTunnelblickProfileArgs `positional-args:"true"`

	Service svc.Service
}

var _ flags.Commander = CreateTunnelblickProfile{}

type createTunnelblickProfileArgs struct {
	DestinationDir string `positional-arg-name:"DESTINATION-DIR" description:"Directory where the *.tblk profile will be created (default: $PWD)"`
}

func (c CreateTunnelblickProfile) Execute(_ []string) error {
	return c.Service.CreateTunnelblickProfile(c.ServiceName, svc.CreateTunnelblickProfileOpts{
		SkipAuthRetry: c.SkipAuthRetry,
		SssocaExec:    c.SssocaExec,
		FileName:      c.Name,
		Directory:     c.Args.DestinationDir,
	})
}
