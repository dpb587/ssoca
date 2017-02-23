# OpenVPN (`openvpn`)

The OpenVPN service is used to sign a user's certificate request and receive a generated connection profile which can be used to connect to a VPN.

This service provides a client command:

* `base-profile` ([details](base-profile-cmd)) - show the base connection profile of the OpenVPN server
* `connect` ([details](connect-cmd)) - connect to a remote OpenVPN server
* `create-profile` ([details](create-profile-cmd)) - create and sign an OpenVPN configuration profile
* `create-tunnelblick-profile` ([details](create-tunnelblick-profile-cmd)) - create a Tunnelblick profile


## Options

 * **`profile`** - the OVPN profile configuration defining the user-agnostic client connection parameters
 * `certauth` - the name of a configured certificate authority (default `default`)
 * `validity` - a [duration](https://golang.org/pkg/time/#ParseDuration) of time for which certificates are signed for (default `2m`)


## External Configuration

OpenVPN is often deployed with CA-based authentication, which is a prerequisite for this service. Ensure the following directives are configured on the OpenVPN server.

    ca_crt ...
    crl_pem ...
