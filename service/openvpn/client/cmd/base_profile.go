package cmd

import (
	clientcmd "github.com/dpb587/ssoca/client/cmd"
	"github.com/jessevdk/go-flags"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type BaseProfile struct {
	clientcmd.ServiceCommand
	clientcmd.InteractiveAuthCommand

	GetClient GetClient
}

var _ flags.Commander = BaseProfile{}

func (c BaseProfile) Execute(_ []string) error {
	client, err := c.GetClient(c.ServiceName, c.SkipAuthRetry)
	if err != nil {
		return bosherr.WrapError(err, "Getting client")
	}

	profile, err := client.BaseProfile()
	if err != nil {
		return bosherr.WrapError(err, "Getting base profile")
	}

	c.Runtime.GetUI().PrintBlock(profile)

	return nil
}
