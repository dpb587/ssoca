// Provide signing services for SSH servers trusting a specific CA.
package ssh

type ServiceType struct{}

func (ServiceType) Type() string {
	return "ssh"
}

func (ServiceType) Version() string {
	return "0.1.0"
}
