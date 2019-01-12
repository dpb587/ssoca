package client

import (
	"fmt"

	"github.com/dpb587/go-onc/onc"
	"github.com/dpb587/go-openvpn/ovpn"
	ovpnonc "github.com/dpb587/go-openvpn/ovpn/onc"
	"github.com/pkg/errors"
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
		return nil, errors.Wrap(err, "Parsing profile")
	}

	oncEncoded, err := ovpnonc.Encode(parsedProfile)
	if err != nil {
		return nil, errors.Wrap(err, "Building ONC")
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
