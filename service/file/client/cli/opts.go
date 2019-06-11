package cli

import (
	"fmt"

	clientcmd "github.com/dpb587/ssoca/client/cmd"
	svc "github.com/dpb587/ssoca/service/file/client"
)

type ServiceCommand struct {
	clientcmd.ServiceCommand
}

func (sc ServiceCommand) GetService() (*svc.Service, error) {
	srv, err := sc.ServiceCommand.GetService()
	if err != nil {
		return nil, err
	}

	mysrv, ok := srv.(*svc.Service)
	if !ok {
		return nil, fmt.Errorf("expected service of type %s but got %s", svc.Service{}.Type(), srv.Type())
	}

	return mysrv, nil
}
