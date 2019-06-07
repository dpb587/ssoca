package openvpn

type ServiceType struct{}

func (ServiceType) Type() string {
	return "openvpn"
}

func (ServiceType) Version() string {
	return "0.1.0"
}
