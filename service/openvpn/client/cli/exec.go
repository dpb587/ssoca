package cli

import (
	"fmt"

	clientcmd "github.com/dpb587/ssoca/client/cmd"
	svc "github.com/dpb587/ssoca/service/openvpn/client"
	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
)

type Exec struct {
	*clientcmd.ServiceCommand `no-flag:"true"`
	clientcmd.InteractiveAuthCommand

	Exec           string      `long:"exec" description:"Path to the openvpn binary"`
	Reconnect      bool        `long:"reconnect" description:"Reconnect on connection disconnects"`
	ManagementMode string      `long:"management-mode" description:"Configure use of management interface (one of: auto, enabled, disabled; default: auto)"`
	Sudo           bool        `long:"sudo" description:"Execute openvpn with sudo"`
	Args           connectArgs `positional-args:"true"`

	// deprecated
	StaticCertificate bool `long:"static-certificate" description:"Write a static certificate in the configuration instead of dynamic renewals (deprecated: use --management-mode=disabled)"`

	serviceFactory svc.ServiceFactory
}

var _ flags.Commander = Exec{}

type connectArgs struct {
	Extra []string `positional-arg-name:"EXTRA" description:"Additional arguments to pass to openvpn"`
}

func (c Exec) Execute(_ []string) error {
	service := c.serviceFactory.New(c.ServiceName)

	if c.StaticCertificate {
		if c.ManagementMode != "" {
			return fmt.Errorf("only one argument may be specified: --management, --static-certificate")
		}

		c.ManagementMode = string(svc.ExecuteManagementModeDisabled)

		c.GetLogger().Warnf("cli: --static-certificate option is deprecated (remove the flag or learn more about the use of --management-mode)")
	} else if c.ManagementMode == "" {
		c.ManagementMode = string(svc.ExecuteManagementModeAuto)
	}

	managementMode, err := svc.ExecuteManagementModeFromString(c.ManagementMode)
	if err != nil {
		return errors.Wrap(err, "parsing --management")
	}

	for {
		err := service.Execute(svc.ExecuteOptions{
			Exec:           c.Exec,
			ExtraArgs:      c.Args.Extra,
			ManagementMode: managementMode,
			Sudo:           c.Sudo,
		})

		if err != nil {
			return err
		} else if !c.Reconnect {
			return nil
		}
	}
}
