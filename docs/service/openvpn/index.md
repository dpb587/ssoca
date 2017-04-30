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


### Renegotiation & Disconnects

By default, OpenVPN attempts to renegotiate the data channel once every hour. Using short-lived certificates will cause the server to terminate the connection when renegotiation occurs (because the certificate has expired). Some clients may attempt to reconnect automatically, however there will typically be a brief interruption in network connectivity.

This behavior is due to OpenVPN not using TLS session renegotiation. Instead, when renegotiating, a full handshake is performed from scratch. Two workarounds are presented below, but you may want to review OpenVPN's technical overview of the [cryptographic layer](https://openvpn.net/index.php/open-source/documentation/security-overview.html) before making a decision.


#### Extended Certificate Lifetime

To retain the default renegotiation behavior of OpenVPN (recommended), you can increase the lifetime that certificates are signed for. For example, setting `validity` to `24h` would allow renegotiations for a day before the connection would fail.

**Note**: this goes against a core principal of `ssoca` which heavily promotes short-lived tokens - a certificate signed at 09:00 could still be used to connect at 21:00. However, the risk of these extended certificates can be reduced by configuring the OpenVPN server with additional verification checks.


##### Extended Certificate Verification

OpenVPN supports a `tls-verify {cmd}` directive which executes external command `{cmd}` to perform final verifications of a peer before it becomes trusted. A script could be used to verify that initial connections occur within 2 minutes from when the certificate was issued (emulating a short-lived token). When renegotiations occur, a script can check whether the client is already trusted and skip the validity checks. An example script is available in the [`ssoca-openvpn-verify`](https://github.com/dpb587/ssoca-bosh-release/tree/src/ssoca-openvpn-verify.go) BOSH job, and OpenVPN could be configured to use it with the following.

    script-security 2
    tls-verify "/var/vcap/packages/ssoca-openvpn-verify/bin/tls-verify 2m"
    tls-export-cert /var/vcap/data/ssoca-openvpn-verify/certs


#### Disable / Increase Renegotiation Timeframe

The renegotiation time is configured with the `reneg-sec {s}` directive where `{s}` is number of seconds and the default is `3600`. This can be increased (for example, `86400` to attempt and renegotiate once per day), which will delay the server from realizing a short-lived certificate was used. Alternatively, the time-based renegotiation can be disabled by setting the value to `0`.

**Note**: disabling or increasing the renegotiation time can theoretically impact the security of your VPN connection.
