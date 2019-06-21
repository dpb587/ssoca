---
title: Cloud Foundry UAA Authentication (uaa-auth)
aliases:
- /auth/authn/uaa/
- /reference/auth/authn/uaa/
---

# Cloud Foundry UAA Authentication (`uaa-auth`)

Cloud Foundry UAA authenticates users from an external UAA server. Users authenticate via CLI prompts or are redirected to UAA and given a token for their CLI usage.


## Server Configuration Options

 * **`url`** -- the address of the UAA server
 * **`public_key`** -- a PEM-formatted public key for verifying JWT tokens
 * **`client_id`** -- the Client ID for authenticating users
 * `client_secret` -- the Client secret for authenticating users
 * `ca_certificate` -- a PEM-formatted certificate for trusting HTTPS connections to the UAA server
 * `prompts` -- a list of prompts to show the user when they are authenticating

{{< note type="warning" >}}
  Client ID and Secret are provided to `ssoca-client` in order for it to communicate with UAA during authentication. Because these tokens are available prior to authentication, they should not be considered secret.
{{< /note >}}


## Authentication Scopes

All scopes propagated by the UAA server will be available in the user's authentication token.


## UAA Client Configuration

In order for `ssoca-client` to connect to UAA, you will need to configure a UAA client for it to use. The following configures [`uaa`](https://bosh.io/jobs/uaa?source=github.com/cloudfoundry/uaa-release&version=67.0#p%3duaa.clients) with a `ssoca_client` ID (no secret), which is allowed to propagate scopes named `env.*`.

```yaml
uaa:
  clients:
    ssoca_client:
      override: true
      authorized-grant-types: "password,refresh_token"
      scope: "openid,env.*"
      authorities: "uaa.none"
      access-token-validity: 120 # 2 min
      refresh-token-validity: 86400 # 1 day
      secret: ""
```
