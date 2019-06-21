---
title: Authorization
aliases:
- /auth/authz/
---

# Authorization

Global (`server.require`) and service-specific (`services[*].require`) authorization restrictions may be configured. For both settings, the values should be an array of hashes with one of the available filter types. At least one authorization rule must be effective for a service.


## `and`

Require multiple filters to be true. The array values of the node must include one or more additional filters.

```yaml
- and:
  - scope: { present: "acme/prod-team" }
  - scope: { present: "acme/security" }
```


## `authenticated`

Require that the user has been successfully authenticated. This filter has no options.

```yaml
- authenticated: ~
```


## `or`

Require at least one of multiple filters to be true. The array values of the node must include one or more additional filters.

```yaml
- or:
  - remote_ip: { within: "192.0.2.0/24" }
  - remote_ip: { within: "198.51.100.0/24" }
```


## `public`

Allow all access, authenticated or anonymous.

```yaml
- public: ~
```


## `remote_ip`

Require that the current API request has come from a specific IP or CIDR.

```yaml
- remote_ip: { within: "192.0.2.1" }
- remote_ip: { within: "192.0.2.1/24" }
- remote_ip: { within: "::1/128" }
```


## `scope`

Require the current authenticated user to have a specific scope (implies `authenticated`).

```yaml
- scope: { present: acme/prod-team }
```


## `service`

Require the current authenticated user to have authenticated with a specific service name (implies `authenticated`).

```yaml
- service: { is: google-cloud }
```


## `username`

Require the current authenticated user to have a specific username (implies `authenticated`).

```yaml
- username: { is: dpb587 }
```
