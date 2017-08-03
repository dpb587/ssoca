package cmd

import (
	clientcmd "github.com/dpb587/ssoca/client/cmd"
	"github.com/dpb587/ssoca/service/openvpn/client/profile"
	"github.com/jessevdk/go-flags"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type CreateProfile struct {
	clientcmd.ServiceCommand
	clientcmd.InteractiveAuthCommand

	GetClient GetClient
}

var _ flags.Commander = CreateProfile{}

func (c CreateProfile) Execute(_ []string) error {
	client, err := c.GetClient(c.ServiceName, c.SkipAuthRetry)
	if err != nil {
		return bosherr.WrapError(err, "Getting client")
	}

	profileManager, err := profile.CreateManagerAndPrivateKey(client, c.ServiceName)
	if err != nil {
		return bosherr.WrapError(err, "Getting profile manager")
	}

	profile, err := profileManager.GetProfile()
	if err != nil {
		return bosherr.WrapError(err, "Getting profile")
	}

	c.Runtime.GetUI().PrintBlock(profile.FullConfig())

	return nil
}
