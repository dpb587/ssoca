// certauth helps manage usages and invocations of Certificate Authority related
// tasks. Different backends may be implemented by providers. For example, a
// filesystem-backed Certificate Authority might use keys from a file, whereas a
// CA backed by Key Management Services in Amazon Web Services may dynamically
// load certificate information as needed.
package certauth

import (
	"crypto/x509"

	"github.com/sirupsen/logrus"

	"golang.org/x/crypto/ssh"
)

const DefaultName = "default"

// Factory allows multiple certificate authority backends to be used.
//
//go:generate counterfeiter . Factory
type Factory interface {
	Create(string, string, map[string]interface{}) (Provider, error)
}

// Manager allows multiple certificate authorities to be referenced.
//
//go:generate counterfeiter . Manager
type Manager interface {
	Add(Provider)
	Get(string) (Provider, error)
}

// Provider is a certificate authority which can access private keys and allow
// them to be used by dependents.
//
//go:generate counterfeiter . Provider
type Provider interface {
	Name() string
	GetCertificate() (*x509.Certificate, error)
	GetCertificatePEM() (string, error)
	SignCertificate(*x509.Certificate, interface{}, logrus.Fields) ([]byte, error)
	SignSSHCertificate(*ssh.Certificate, logrus.Fields) error
}

// ProviderFactory handles the creation of a specific certificate authority
// provider with its configuration.
//
//go:generate counterfeiter . ProviderFactory
type ProviderFactory interface {
	Create(string, map[string]interface{}) (Provider, error)
}
