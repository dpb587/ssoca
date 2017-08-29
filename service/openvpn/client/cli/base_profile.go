package cli

import (
	clientcmd "github.com/dpb587/ssoca/client/cmd"
	svc "github.com/dpb587/ssoca/service/openvpn/client"
	"github.com/jessevdk/go-flags"
)

type BaseProfile struct {
	clientcmd.ServiceCommand
	clientcmd.InteractiveAuthCommand

	serviceFactory svc.ServiceFactory
}

var _ flags.Commander = BaseProfile{}

func (c BaseProfile) Execute(_ []string) error {
	service := c.serviceFactory.New(c.ServiceName)

	profile, err := service.BaseProfile(svc.BaseProfileOptions{
		SkipAuthRetry: c.SkipAuthRetry,
	})
	if err != nil {
		return err
	}

	c.Runtime.GetUI().PrintBlock(profile)

	return nil
}
