package openvpn

type Service struct{}

func (Service) Type() string {
	return "openvpn"
}

func (Service) Version() string {
	return "0.1.0"
}
