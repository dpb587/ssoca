# Authorization

Global (`auth.require`) and service-specific (`services[*].require`) authorization restrictions may be configured. For both settings, the values should be an array of hashes with one of the available filter types.

By default, services are accessible.


## `authenticated`

Require that the user has been successfully authenticated. This filter has no options.

    - authenticated: ~


## `remote_ip`

Require that the current API request has come from a specific IP or CIDR.

    - remote_ip: { within: "192.0.2.1" }
    - remote_ip: { within: "192.0.2.1/24" }
    - remote_ip: { within: "::1/128" }


## `scope`

Require that a specific scope is present for the current authenticated user (implies `authenticated`).

    - scope: { present: acme/prod-team }
