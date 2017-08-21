package cli

import (
	clientcmd "github.com/dpb587/ssoca/client/cmd"
	svc "github.com/dpb587/ssoca/service/openvpn/client"
	"github.com/jessevdk/go-flags"
)

type BaseProfile struct {
	clientcmd.ServiceCommand
	clientcmd.InteractiveAuthCommand

	Service svc.Service
}

var _ flags.Commander = BaseProfile{}

func (c BaseProfile) Execute(_ []string) error {
	profile, err := c.Service.BaseProfile(c.ServiceName, svc.BaseProfileOptions{
		SkipAuthRetry: c.SkipAuthRetry,
	})
	if err != nil {
		return err
	}

	c.Runtime.GetUI().PrintBlock(profile)

	return nil
}
