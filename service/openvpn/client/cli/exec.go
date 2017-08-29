package cli

import (
	clientcmd "github.com/dpb587/ssoca/client/cmd"
	svc "github.com/dpb587/ssoca/service/openvpn/client"
	"github.com/jessevdk/go-flags"
)

type Exec struct {
	clientcmd.ServiceCommand
	clientcmd.InteractiveAuthCommand

	Exec              string      `long:"exec" description:"Path to the openvpn binary"`
	Reconnect         bool        `long:"reconnect" description:"Reconnect on connection disconnects"`
	StaticCertificate bool        `long:"static-certificate" description:"Write a static certificate in the configuration instead of dynamic renewals"`
	Sudo              bool        `long:"sudo" description:"Execute openvpn with sudo"`
	Args              connectArgs `positional-args:"true"`

	serviceFactory svc.ServiceFactory
}

var _ flags.Commander = Exec{}

type connectArgs struct {
	Extra []string `positional-arg-name:"EXTRA" description:"Additional arguments to pass to openvpn"`
}

func (c Exec) Execute(_ []string) error {
	service := c.serviceFactory.New(c.ServiceName)

	for {
		err := service.Execute(svc.ExecuteOptions{
			StaticCertificate: c.StaticCertificate,
			Sudo:              c.Sudo,
			Exec:              c.Exec,
			ExtraArgs:         c.Args.Extra,
		})

		if err != nil {
			return err
		} else if !c.Reconnect {
			return nil
		}
	}
}
