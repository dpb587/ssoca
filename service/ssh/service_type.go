package ssh

import "github.com/dpb587/ssoca/service"

const Type service.Type = "ssh"

type ServiceType struct{}

func (ServiceType) Type() service.Type {
	return Type
}

func (ServiceType) Version() string {
	return "0.1.0"
}
