# Google Authentication (`google`)

Google authenticates users through their Google accounts. Users are redirected through Google's website and given a token for their CLI usage. Once authenticated, short-lived [JSON Web Tokens](../jwt) are used to persist and validate authentication state.


## Options

 * **`client_id`** - Client ID from registered application
 * **`client_secret`** - Client Secret from registered application
 * **`jwt`** - a hash of JWT signing details
    * **`private_key`** - a PEM-formatted private key
    * `validity` - a [duration](https://golang.org/pkg/time/#ParseDuration) for how long authentication tokens will be remembered (default `24h`)
 * `auth_url` - authentication URL (default `https://accounts.google.com/o/oauth2/v2/auth`)
 * `token_url` - token URL (default `https://www.googleapis.com/oauth2/v4/token`)
 * `scopes` - optionally load additional profile information for extended scopes
    * `cloud_project` - request information from Google Cloud (requires the [Google Cloud Resource Manager API](https://console.cloud.google.com/apis/api/cloudresourcemanager.googleapis.com/overview) to be enabled)
       * `projects` - a list of project identifiers to check for membership; if left empty, all projects will be checked
       * `roles` - a list of roles (e.g. `roles/owner`) to check for access; if left empty, all roles will be included


## Authentication Scopes

When a user authenticates, their email, email mailbox, and email domain will be added as scopes.

 * `somebody@example.com`
 * `email/mailbox/somebody`
 * `email/domain/example.com`

If `cloud_project` scopes are enabled, their project role scopes will also be added as scopes.

 * `cloud/project/1234567890/roles/owner`


## Google Application

This requires [registering a credential](https://console.cloud.google.com/apis/credentials). At a minimum, ensure the following fields are configured.

 * Credential type &ndash; OAuth client ID
 * Application type &ndash; Web application
 * Authorized redirect URIs &ndash; `https://{ssoca_host}:{ssoca_port}/auth/callback`


## General Notes

 * changing the `jwt.private_key` will revoke all existing authentication tokens
