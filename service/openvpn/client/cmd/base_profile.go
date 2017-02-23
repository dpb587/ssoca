package cmd

import (
	clientcmd "github.com/dpb587/ssoca/client/cmd"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type BaseProfile struct {
	clientcmd.ServiceCommand

	GetClient GetClient
}

func (c *BaseProfile) Execute(args []string) error {
	client, err := c.GetClient(c.ServiceName)
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
