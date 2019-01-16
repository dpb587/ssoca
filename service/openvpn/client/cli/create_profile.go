package cli

import (
	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"

	clientcmd "github.com/dpb587/ssoca/client/cmd"
	svc "github.com/dpb587/ssoca/service/openvpn/client"
)

type CreateProfile struct {
	*clientcmd.ServiceCommand `no-flag:"true"`
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
		return errors.Wrap(err, "Getting profile")
	}

	c.Runtime.GetUI().PrintBlock([]byte(profile.StaticConfig()))

	return nil
}
