# GitHub (`github`)

GitHub authenticates users through their GitHub accounts. Users are redirected through GitHub's website and given a token for their CLI usage. Once authenticated, short-lived [JSON Web Tokens](../jwt) are used to persist and validate authentication state.


## Options

 * **`client_id`** - Client ID from registered application
 * **`client_secret`** - Client Secret from registered application
 * **`jwt`** - a hash of JWT signing details
    * **`private_key`** - a PEM-formatted private key
    * `validity` - a [duration](https://golang.org/pkg/time/#ParseDuration) for how long authentication tokens will be remembered (default `24h`)
 * `auth_url` - authentication URL (default `https://github.com/login/oauth/authorize`)
 * `token_url` - token URL (default `https://github.com/login/oauth/access_token`)


## Authentication Scopes

When a user authenticates, their organization and team membership information will be pulled and converted into scopes.

Examples

 * `dpb587` (user)
 * `theloopyewe` (organization membership)
 * `cloudfoundry/open-source-contributor` (organization team membership)

When a user's organizations or teams change, they will need to logout and log back in before their scopes will be updated.


## GitHub Application

This requires [registering an application](https://github.com/settings/applications/new). At a minimum, ensure the following fields are configured.

 * Authorization Callback URL &ndash; `https://{ssoca_host}:{ssoca_port}/auth/callback`


## General Notes

 * changing the `jwt.private_key` will revoke all existing authentication tokens
