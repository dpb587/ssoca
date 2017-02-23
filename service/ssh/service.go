// Provide signing services for SSH servers trusting a specific CA.
package ssh

type Service struct{}

func (Service) Type() string {
	return "ssh"
}

func (Service) Version() string {
	return "0.1.0"
}
