package onc

import (
	"encoding/base64"
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/dpb587/go-onc/onc"
	"github.com/dpb587/go-onc/onc/oncutil"
	"github.com/dpb587/go-openvpn/ovpn"
)

var whitespaceRE = regexp.MustCompile(`\s+`)

func Encode(profile *ovpn.Profile) (*onc.ONC, error) {
	openvpnConfig := onc.OpenVPN{}
	networkConfigurationConfig := onc.NetworkConfiguration{
		GUID: "openvpn",
		Name: profile.Name,
		Type: "VPN",
		VPN: onc.VPN{
			Type: "OpenVPN",
		},
	}
	config := onc.ONC{
		Type: "UnencryptedConfiguration",
	}

	var certBytes, keyBytes []byte

	for _, e := range profile.Elements {
		switch pe := e.(type) {
		case ovpn.DirectiveProfileElement:
			switch pe.Directive() {
			case "auth":
				openvpnConfig.Auth = pe.Args()[0]
			case "auth-retry":
				openvpnConfig.AuthRetry = pe.Args()[0]
			case "auth-nocache":
				openvpnConfig.AuthNoCache = true
			case "cipher":
				openvpnConfig.Cipher = pe.Args()[0]
			case "comp-lzo":
				v := pe.Args()[0]
				switch v {
				case "yes":
					v = "true"
				case "no":
					v = "false"
				}

				openvpnConfig.CompLZO = v
			case "comp-noadapt":
				openvpnConfig.CompNoAdapt = true
			case "key-direction":
				openvpnConfig.KeyDirection = pe.Args()[0]
			case "ns-cert-type":
				openvpnConfig.NsCertType = pe.Args()[0]
			case "push-peer-info":
				openvpnConfig.PushPeerInfo = true
			case "proto":
				openvpnConfig.Proto = pe.Args()[0]
			case "remote":
				subpieced := whitespaceRE.Split(pe.Args()[0], 3)
				if networkConfigurationConfig.VPN.Host == "" {
					networkConfigurationConfig.VPN.Host = pe.Args()[0]

					if len(pe.Args()) > 1 {
						port, err := strconv.Atoi(pe.Args()[1])
						if err != nil {
							return nil, fmt.Errorf("invalid remote port: %s", subpieced[1])
						}

						openvpnConfig.Port = port

						if len(pe.Args()) > 2 {
							openvpnConfig.Proto = pe.Args()[2]
						}
					}
				} else {
					openvpnConfig.ExtraHosts = append(openvpnConfig.ExtraHosts, subpieced[0])

					// TODO validate matching port/proto?
				}
			case "remote-cert-eku":
				openvpnConfig.RemoteCertEKU = pe.Args()[0]
			case "remote-cert-ku":
				openvpnConfig.RemoteCertKU = append(openvpnConfig.RemoteCertKU, pe.Args()[0])
			case "remote-cert-tls":
				openvpnConfig.RemoteCertTLS = pe.Args()[0]
			case "reneg-sec":
				renegSec, err := strconv.Atoi(pe.Args()[0])
				if err != nil {
					return nil, fmt.Errorf("invalid reneg-sec: %s", pe.Args()[0])
				}

				openvpnConfig.RenegSec = renegSec
			case "server-poll-timeout", "connect-timeout":
				serverPollTimeout, err := strconv.Atoi(pe.Args()[0])
				if err != nil {
					return nil, fmt.Errorf("invalid server-poll-timeout: %s", pe.Args()[0])
				}

				openvpnConfig.ServerPollTimeout = serverPollTimeout
			case "shaper":
				shaper, err := strconv.Atoi(pe.Args()[0])
				if err != nil {
					return nil, fmt.Errorf("invalid shaper: %s", pe.Args()[0])
					panic("invalid shaper")
				}

				openvpnConfig.Shaper = shaper
			case "static-challenge":
				subpieced := whitespaceRE.Split(pe.Args()[0], 2)
				if len(subpieced) != 2 {
					return nil, fmt.Errorf("invalid static-challenge: %s", pe.Args()[0])
				} else if subpieced[1] != "1" {
					return nil, fmt.Errorf("onc only support echoing static-challenge: %s", pe.Args()[0])
				}

				openvpnConfig.StaticChallenge = subpieced[0]
			case "tls-version-min":
				openvpnConfig.TLSVersionMin = pe.Args()[0]
			case "verb":
				openvpnConfig.Verb = pe.Args()[0]
			case "verify-hash":
				subpieced := whitespaceRE.Split(pe.Args()[0], 2)
				if len(subpieced) == 2 && subpieced[1] != "SHA1" {
					return nil, fmt.Errorf("onc only supports verify-hash with sha1")
				}

				openvpnConfig.VerifyHash = subpieced[0]
			case "verify-x509-name": // TODO 'quoted name' // TODO unsupported type
				openvpnConfig.VerifyX509 = pe.Args()[0]
			}
		case ovpn.EmbeddedProfileElement:
			switch pe.Embed() {
			case "ca":
				config.Certificates = append(
					config.Certificates,
					onc.Certificate{
						GUID: "ca",
						Type: "Server",
						X509: pe.Data(),
					},
				)

				openvpnConfig.ServerCARefs = append(openvpnConfig.ServerCARefs, "ca")
			case "cert":
				certBytes = []byte(pe.Data())
			case "key":
				keyBytes = []byte(pe.Data())
			case "tls-auth":
				openvpnConfig.TLSAuthContents = pe.Data()
			}
		}
	}

	if len(certBytes) > 0 || len(keyBytes) > 0 {
		pkcsBytes, err := oncutil.ConvertKeyPairToPKCS12(certBytes, keyBytes)
		if err != nil {
			log.Printf("ERROR: converting to pkcs12: %v", err)
			panic("converting to pkcs12")
		}

		config.Certificates = append(
			config.Certificates,
			onc.Certificate{
				GUID:   "client",
				PKCS12: base64.StdEncoding.EncodeToString(pkcsBytes),
				Type:   "Client",
			},
		)

		openvpnConfig.ClientCertType = "Ref"
		openvpnConfig.ClientCertRef = "client"
	}

	networkConfigurationConfig.VPN.OpenVPN = openvpnConfig
	config.NetworkConfigurations = append(
		config.NetworkConfigurations,
		networkConfigurationConfig,
	)

	return &config, nil
}
