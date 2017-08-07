package env

type Service struct{}

func (Service) Type() string {
	return "env"
}

func (Service) Version() string {
	return "0.6.0"
}
