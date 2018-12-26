---
title: OpenVPN (openvpn)
aliases:
- /service/openvpn/
---

# OpenVPN (`openvpn`)

The OpenVPN service is used to sign a user's certificate request and receive a generated connection profile which can be used to connect to a VPN.

Before a client can use this service, they must install OpenVPN ([learn more](https://github.com/dpb587/openvpn-bosh-release/blob/master/docs/ops/users/software.md)).


## Client Commands

The OpenVPN workflows can be managed through `openvpn` subcommands.

* `base-profile` ([details]({{< ref "/reference/service/openvpn/base-profile-cmd.md" >}})) - show the base connection profile of the OpenVPN server
* `create-launchd-service` ([details]({{< ref "/reference/service/openvpn/create-launchd-service-cmd.md" >}})) - create a launchd service
* `create-onc-profile` ([details]({{< ref "/reference/service/openvpn/create-onc-profile-cmd.md" >}})) - create an ONC profile
* `create-profile` ([details]({{< ref "/reference/service/openvpn/create-profile-cmd.md" >}})) - create and sign an OpenVPN configuration profile
* `create-tunnelblick-profile` ([details]({{< ref "/reference/service/openvpn/create-tunnelblick-profile-cmd.md" >}})) - create a Tunnelblick profile
* `exec` ([details]({{< ref "/reference/service/openvpn/exec-cmd.md" >}})) - connect to a remote OpenVPN server


## Server Configuration Options

The following may be configured in the `options` section when configuring an `openvpn` [service]({{< ref "/reference/server/configuration.md#services" >}}).

 * **`profile`** - the OVPN profile configuration defining the user-agnostic client connection parameters
 * `certauth` - the name of a configured certificate authority (default `default`)
 * `validity` - a [duration](https://golang.org/pkg/time/#ParseDuration) of time for which certificates are signed for (default `2m`)

*Tip*: the OpenVPN server must be configured to trust certificates signed by ssoca and to optionally enforce extended certificate validity ([learn more]({{< ref "/reference/service/openvpn/external-configuration.md" >}})).
