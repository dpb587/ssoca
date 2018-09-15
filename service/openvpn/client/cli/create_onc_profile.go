package cli

import (
	"encoding/json"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	clientcmd "github.com/dpb587/ssoca/client/cmd"
	svc "github.com/dpb587/ssoca/service/openvpn/client"
	"github.com/jessevdk/go-flags"
)

type CreateONCProfile struct {
	clientcmd.ServiceCommand
	clientcmd.InteractiveAuthCommand

	serviceFactory svc.ServiceFactory

	Name string `long:"name" description:"Specific network configuration name to use"`
}

var _ flags.Commander = CreateProfile{}

func (c CreateONCProfile) Execute(_ []string) error {
	service := c.serviceFactory.New(c.ServiceName)

	onc, err := service.CreateONCProfile(svc.CreateONCProfileOptions{
		SkipAuthRetry: c.SkipAuthRetry,
		Name:          c.Name,
	})
	if err != nil {
		return bosherr.WrapError(err, "Getting profile")
	}

	oncBytes, err := json.MarshalIndent(onc, "", "  ")
	if err != nil {
		return bosherr.WrapError(err, "Encoding ONC")
	}

	oncBytes = append(oncBytes, '\n')

	c.Runtime.GetUI().PrintBlock(oncBytes)

	return nil
}
