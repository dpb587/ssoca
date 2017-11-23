package cli

import (
	clientcmd "github.com/dpb587/ssoca/client/cmd"
	svc "github.com/dpb587/ssoca/service/openvpn/client"
	"github.com/jessevdk/go-flags"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type CreateProfile struct {
	clientcmd.ServiceCommand
	clientcmd.InteractiveAuthCommand

	serviceFactory svc.ServiceFactory
}

var _ flags.Commander = CreateProfile{}

func (c CreateProfile) Execute(_ []string) error {
	service := c.serviceFactory.New(c.ServiceName)

	profile, err := service.CreateProfile(svc.CreateProfileOptions{
		SkipAuthRetry: c.SkipAuthRetry,
	})
	if err != nil {
		return bosherr.WrapError(err, "Getting profile")
	}

	c.Runtime.GetUI().PrintBlock([]byte(profile.StaticConfig()))

	return nil
}
