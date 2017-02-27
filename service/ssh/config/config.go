package config

import (
	"time"

	"github.com/dpb587/ssoca/certauth"
	"github.com/dpb587/ssoca/server/service/dynamicvalue"
)

type configCriticalOption string
type configExtension string

const (
	// CriticalOptionForceCommand defines a command that is executed (replacing any the user specified on the ssh command-line) whenever this key is used for authentication.
	CriticalOptionForceCommand configCriticalOption = "force-command"

	// CriticalOptionSourceAddress defines a comma-separated list of source addresses from which this certificate is accepted for authentication. Addresses are specified in CIDR format (nn.nn.nn.nn/nn or hhhh::hhhh/nn). If this option is not present then certificates may be presented from any source address.
	CriticalOptionSourceAddress configCriticalOption = "source-address"

	// ExtensionPermitX11Forwarding is a flag indicating that X11 forwarding should be permitted. X11 forwarding will be refused if this option is absent.
	ExtensionPermitX11Forwarding configExtension = "permit-X11-forwarding"

	// ExtensionPermitAgentForwarding is a flag indicating that agent forwarding should be allowed. Agent forwarding must not be permitted unless this option is present.
	ExtensionPermitAgentForwarding configExtension = "permit-agent-forwarding"

	// ExtensionPermitPortForwarding is a flag indicating that port-forwarding should be allowed. If this option is not present then no port forwarding will be allowed.
	ExtensionPermitPortForwarding configExtension = "permit-port-forwarding"

	// ExtensionPermitPTY is a flag indicating that PTY allocation should be permitted. In the absence of this option PTY allocation will be disabled.
	ExtensionPermitPTY configExtension = "permit-pty"

	// ExtensionPermitUserRC is a flag indicating that execution of ~/.ssh/rc should be permitted. Execution of this script will not be permitted if this option is not present.
	ExtensionPermitUserRC configExtension = "permit-user-rc"

	// ExtensionNoDefaults disables the default set of extensions ssoca would normally add.
	ExtensionNoDefaults configExtension = "ssoca-no-defaults"
)

// ExtensionDefaults is the set of extensions which will be enabled if ExtensionNoDefaults is not configured.
var ExtensionDefaults = Extensions{
	ExtensionPermitX11Forwarding,
	ExtensionPermitAgentForwarding,
	ExtensionPermitPortForwarding,
	ExtensionPermitPTY,
	ExtensionPermitUserRC,
}

// Config settings for SSH key signing.
type Config struct {
	CertAuthName   string   `yaml:"certauth,omitempty"`
	ValidityString string   `yaml:"validity,omitempty"`
	RawPrincipals  []string `yaml:"principals,omitempty"`

	Target Target `yaml:"target,omitempty"`

	CriticalOptions CriticalOptions `yaml:"critical_options,omitempty"`
	Extensions      Extensions      `yaml:"extensions,omitempty"`

	CertAuth   certauth.Provider    `yaml:"-"`
	Validity   time.Duration        `yaml:"-"`
	Principals []dynamicvalue.Value `yaml:"-"`
}

type CriticalOptions map[configCriticalOption]string

type Extensions []configExtension

type Target struct {
	Host    string `yaml:"host,omitempty" json:"host,omitempty"`
	RawUser string `yaml:"user,omitempty" json:"user,omitempty"`
	Port    int    `yaml:"port,omitempty" json:"port,omitempty"`

	User dynamicvalue.Value `yaml:"-"`
}

func (c Target) Configured() bool {
	return c.Host != "" || c.RawUser != "" || c.Port != 0
}

// ApplyDefaults provides some static default values.
func (c *Config) ApplyDefaults() {
	if c.CertAuthName == "" {
		c.CertAuthName = "default"
	}

	if c.ValidityString == "" {
		c.ValidityString = "2m"
	}

	if len(c.Extensions) == 0 {
		c.Extensions = ExtensionDefaults
	}
}
