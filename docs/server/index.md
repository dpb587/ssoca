# Server

The server is configured through a YAML file with several sections. The first section is `server` which configures the HTTPS listener...

    server:
      # tls configuration
      certificate_path: ~ # PEM-formatted certificate (required)
      private_key_path: ~ # PEM-formatted private key (required)

      # bind settings
      host: 0.0.0.0 # default
      port: 18705   # default

      # optional redirect of /
      root_redirect: ~ # optional

The YAML configuration should also include an `auth` section referencing a [service type](../authn) and the services options...

    auth:
      type: uaa
      options:
        public_key: ...snip...

Certificate authorities can be defined in the `certauths` field which is an array of CA providers referencing a [CA type](../certauth) and the CA options. Services may later reference CAs by their name...

    certauths:
      - type: fs # one of the available CA types
        name: ~  # defaults to `default`
        options: # CA-specific options
          private_key_path: /some/path.crt

The last section is `services` which is an array of service configurations referencing a [service type](../service) and the service options...

    services:
      - type: ssh # one of the available service types
        name: ~   # defaults to the value of `type`
        options:  # service-specific options
          host: 192.0.2.1
          user: vcap
