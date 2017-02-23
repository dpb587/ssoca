package helper

import boshuaa "github.com/cloudfoundry/bosh-cli/uaa"

//go:generate counterfeiter . ClientFactory
type ClientFactory interface {
	CreateClient(string, string, string, string) (boshuaa.UAA, error)
}
