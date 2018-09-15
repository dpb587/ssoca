package client

import (
	"fmt"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/dpb587/go-onc/onc"
	"github.com/dpb587/go-openvpn/ovpn"
	ovpnonc "github.com/dpb587/go-openvpn/ovpn/onc"
)

type CreateONCProfileOptions struct {
	SkipAuthRetry bool
	Name          string
}

func (s Service) CreateONCProfile(opts CreateONCProfileOptions) (*onc.ONC, error) {
	profile, err := s.CreateProfile(CreateProfileOptions{
		SkipAuthRetry: opts.SkipAuthRetry,
	})
	if err != nil {
		return nil, err
	}

	parsedProfile, err := ovpn.Parse([]byte(profile.StaticConfig()))
	if err != nil {
		return nil, bosherr.WrapError(err, "Parsing profile")
	}

	oncEncoded, err := ovpnonc.Encode(parsedProfile)
	if err != nil {
		return nil, bosherr.WrapError(err, "Building ONC")
	}

	name := opts.Name
	if name == "" {
		name = s.runtime.GetEnvironmentName()

		if s.name != "openvpn" {
			name = fmt.Sprintf("%s-%s", name, s.name)
		}
	}

	oncEncoded.NetworkConfigurations[0].Name = name
	oncEncoded.NetworkConfigurations[0].VPN.OpenVPN.UserAuthenticationType = "None"

	return oncEncoded, nil
}
