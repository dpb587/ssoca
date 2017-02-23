package auth

type Service struct{}

func (Service) Type() string {
	return "auth"
}

func (Service) Version() string {
	return "0.1.0"
}
