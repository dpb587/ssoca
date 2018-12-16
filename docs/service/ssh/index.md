# SSH (`ssh`)

The SSH service is used to sign a user's public key for accessing an intentionally-configured SSH servers.

This service provides a client command:

* `agent` ([details](agent-cmd)) - start an SSH agent
* `exec` ([details](exec-cmd)) - connect to a remote SSH server
* `sign-public-key` ([details](sign-public-key-cmd)) - create a certificate for a specific public key


## Options

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


## General Notes

* certificates are only used for authentication (sessions are not disconnected once the certificates expire)


## External Configuration

By default, OpenSSH relies on PAM and Public Key authentication (`authorized_keys`). To use this service, you must configure the OpenSSH server to trust users presenting a certificate which has been signed by a particular certificate authority.

First, you should generate the public key of your certificate authority in OpenSSH format. You can use `ssh-keygen` and give it the private key of the certificate authority (this generates a public key and is not secret).

    ssh-keygen -f ca.key -y

Copy the output into a file on the OpenSSH server (e.g. `/etc/ssh/trusted_user_ca_keys`), and update the `TrustedUserCAKeys` setting in `/etc/ssh/sshd_config`.

    TrustedUserCAKeys /etc/ssh/trusted_user_ca_keys

If you do not want to allow users to manage their own `~/.ssh/authorized_keys` file (forcing all public key connections to be signed by the CA), you may want to update the `AuthorizedKeysFile` setting.

    AuthorizedKeysFile /dev/null

Once configured, restart the `ssh` service.

    service ssh restart

For a BOSH-managed server, you may find the [ssh-conf](https://github.com/dpb587/ssh-conf-bosh-release) BOSH release useful.


## Debugging

If a signed certificate is not working, sometimes it's helpful to inspect the signed certificate, taking particular interest in the principals...

    $ ssh-keygen -L -f <( ssoca ssh sign-public-key ~/.ssh/id_rsa.pub )
    /dev/fd/63:
        Type: ssh-rsa-cert-v01@openssh.com user certificate
        Public key: RSA-CERT SHA256:Lbm8fojiin5Mn95obC0Qxxf9/Gca4GtJMuUfax4Vu7M
        Signing CA: RSA SHA256:9cqZE53uBj8fA5MBg9OBU9fzQ6L10G4O90x0ETgFp7E
        Key ID: "somebody@example.com"
        Serial: 0
        Valid: from 2017-02-28T22:53:47 to 2017-02-28T22:55:52
        Principals:
                somebody
                vcap
        Critical Options: (none)
        Extensions:
                permit-X11-forwarding
                permit-agent-forwarding
                permit-port-forwarding
                permit-pty
                permit-user-rc

To convert a X509 private key to an OpenSSH public key...

    $ ssh-keygen -f ca-private.pem -y > ca-cert.pub
