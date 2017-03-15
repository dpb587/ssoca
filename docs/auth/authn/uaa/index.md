# Cloud Foundry UAA (`uaa`)

Cloud Foundry UAA authenticates users from an external UAA server. Users authenticate via CLI prompts or are redirected to UAA and given a token for their CLI usage.


## Options

 * **`url`** - the address of the UAA server
 * **`public_key`** - a PEM-formatted public key for verifying JWT tokens
 * `ca_certificate` - a PEM-formatted certificate for trusting HTTPS connections to the UAA server


## Authentication Scopes

All scopes propagated by the UAA server will be available in the user's authentication token.
