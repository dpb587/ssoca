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

    server:
      # tls configuration
      certificate_path: ~ # PEM-formatted certificate; required
      private_key_path: ~ # PEM-formatted private key; required

      # bind settings
      host: 0.0.0.0 # default
      port: 18705   # default

      # optional redirects
      redirects:
        root: ~         # redirect of /; optional
        auth_failure: ~ # optional
        auth_success: ~ # optional

      # optionally configure proxies you expect to receive traffic from
      trusted_proxies:
      - 127.0.0.1/8
      - ::1

      # optionally configure a robots.txt response (the following is default)
      robotstxt: |
        User-agent: *
        Disallow: /


## `auth`

The YAML configuration must also include an `auth` section referencing an [authentication service type](../auth/authn) and the service's options...

    auth:
      type: uaa
      options:
        public_key: ...snip...


## `certauths`

Certificate authorities can be defined in the `certauths` field which is an array of CA providers referencing a [CA type](../certauth) and the CA options. Services may later reference CAs by their name...

    certauths:
      - type: fs # one of the available CA types
        name: ~  # defaults to `default`
        options: # CA-specific options
          private_key_path: /some/path.crt


## `services`

The last required section is `services` which is an array of service configurations referencing a [service type](../service) and the service options...

    services:
      - type: ssh # one of the available service types
        name: ~   # defaults to the value of `type`
        options:  # service-specific options
          host: 192.0.2.1
          user: vcap


## `env`

Optionally, an `env` section can be configured with some end user-oriented details. Options are documented [here](service/env/#options).
