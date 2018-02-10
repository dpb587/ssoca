package cmd

import (
	"time"

	"github.com/dpb587/ssoca/client"
	"github.com/dpb587/ssoca/version"
	flags "github.com/jessevdk/go-flags"
)

type Version struct {
	Runtime client.Runtime

	Name   bool `long:"name" description:"Show only the application name"`
	Semver bool `long:"semver" description:"Show only the semver version value"`
	Commit bool `long:"commit" description:"Show only the versioning commit reference"`
	Built  bool `long:"built" description:"Show only the build date"`

	Version version.Version
}

var _ flags.Commander = Version{}

func (c Version) Execute(_ []string) error {
	ui := c.Runtime.GetUI()

	if c.Name {
		ui.PrintBlock(append([]byte(c.Version.Name), '\n'))
	} else if c.Semver {
		ui.PrintBlock(append([]byte(c.Version.Semver), '\n'))
	} else if c.Commit {
		ui.PrintBlock(append([]byte(c.Version.Commit), '\n'))
	} else if c.Built {
		ui.PrintBlock(append([]byte(c.Version.Built.Format(time.RFC3339)), '\n'))
	} else {
		ui.PrintLinef(c.Version.String())
	}

	return nil
}
