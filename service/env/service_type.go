package env

import "github.com/dpb587/ssoca/service"

const Type service.Type = "env"

type ServiceType struct{}

func (ServiceType) Type() service.Type {
	return Type
}

func (ServiceType) Version() string {
	return "0.6.0"
}
