package cmd

import (
	clientcmd "github.com/dpb587/ssoca/client/cmd"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type CreateProfile struct {
	clientcmd.ServiceCommand

	CreateUserProfile CreateUserProfile
}

func (c *CreateProfile) Execute(args []string) error {
	profile, err := c.CreateUserProfile(c.ServiceName)
	if err != nil {
		return bosherr.WrapError(err, "Creating profile")
	}

	c.Runtime.GetUI().PrintBlock(profile)

	return nil
}
