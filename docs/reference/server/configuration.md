---
title: Configuration
aliases:
- /server/configuration/
---

# Configuration

The server is configured through a YAML file with several top-level sections.

For a more full example, refer to the [server.conf](https://github.com/dpb587/ssoca-bosh-release/blob/master/jobs/ssoca/templates/etc/server.conf.erb) from the BOSH release.


## `server`

The first section is `server` which configures the HTTPS listener...

```yaml
server:
  # tls configuration
  certificate_path: ~ # PEM-formatted certificate; required
  private_key_path: ~ # PEM-formatted private key; required

  # bind settings
  host: "0.0.0.0" # default
  port: 18705     # default

  # optional redirects
  redirects:
    root: ~         # redirect of /; optional
    auth_failure: ~ # optional
    auth_success: ~ # optional

  # optionally configure upstream proxies (used for remote IP reporting)
  trusted_proxies:
  - "127.0.0.1/8"
  - "::1"

  # optionally configure a robots.txt response (the following is default)
  robotstxt: |
    User-agent: *
    Disallow: /
```


## `certauths`

Certificate authorities can be defined in the `certauths` field which is an array of CA providers referencing a [CA type]({{< relref "../certauth" >}}) and the CA options. Services may later reference CAs by their name...

```yaml
certauths:
  - type: "fs" # one of the available CA types
    name: ~    # defaults to `default`
    options:   # CA-specific options
      private_key_path: "/some/path.crt"
```


## `services`

The last section is `services` which is an array of service configurations referencing a service type and the service options. You will typically configure at least one authentication server and one user service...

```yaml
services:
  - type: "ssh" # one of the available service types
    name: ~     # defaults to the value of `type`
    options:    # service-specific options
      host: "192.0.2.1"
      user: "vcap"
  - type: "github-auth"
    name: "auth"
    options: # ...
```


## `env`

Optionally, an `env` section can be configured with some end user-oriented details. Options are documented [here]({{< relref "../service/env/_index.md#server-configuration-options" >}}).
