---
title: SSH (ssh)
aliases:
- /service/ssh/
---

# SSH (`ssh`)

The SSH service is used to sign a user's public key for accessing an intentionally-configured SSH servers.


## Client Commands

The SSH workflows can be managed through `ssh` subcommands.

* `agent` ([details]({{< ref "/reference/service/ssh/agent-cmd" >}})) - start an SSH agent
* `exec` ([details]({{< ref "/reference/service/ssh/exec-cmd" >}})) - connect to a remote SSH server
* `sign-public-key` ([details]({{< ref "/reference/service/ssh/sign-public-key-cmd" >}})) - create a certificate for a specific public key


## Server Configuration Options

The following may be configured in the `options` section when configuring an `ssh` [service]({{< ref "/reference/server/configuration.md#services" >}}).

* **`principals`** - an array of usernames to allow SSH sessions for ([templatized](../../server/templating))
* `certauth` - the name of a configured certificate authority (default `default`)
* `validity` - a [duration](https://golang.org/pkg/time/#ParseDuration) of time for which certificates are signed for (default `2m`)
* `critical_options` - a *hash* of specific settings further restricting connections to the SSH server
  * `force-command` - a command which is forcefully executed on the SSH server ([templatized](../../server/templating))
  * `source-address` - a CSV list of source addresses in CIDR format which certificates can come from for authentication ([templatized](../../server/templating))
* `extensions` - an *array* of session features for the server to enforce on the connection (default all)
  * `permit-X11-forwarding`
  * `permit-agent-forwarding`
  * `permit-port-forwarding`
  * `permit-pty`
  * `permit-user-rc`

**Propagated Client Options**

* `client` - a hash of settings influencing client behavior
  * `host` - the remote host of the SSH server
  * `port` - the remote port of the SSH server (default `22`)
  * `user` - the remote user to authenticate as ([templatized](../../server/templating))
  * `public_key` - the public key of the remote SSH server (requires `host`)

*Tip*: the SSH server must be configured to trust certificates signed by ssoca ([learn more]({{< ref "/reference/service/ssh/external-configuration.md" >}})).


## General Notes

* certificates are only used for authentication (sessions are not disconnected once the certificates expire)


## Known Issues

 * A recent OS X update introduces a new version of `ssh` which does not correctly check for supported algorithms. Using `brew install openssh` to install a newer version seems to fix the issue ([related](https://bugs.launchpad.net/ubuntu/+source/openssh/+bug/1790963) [bug](http://bugzilla.mindrot.org/show_bug.cgi?id=2799)).
