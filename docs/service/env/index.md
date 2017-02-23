# Environment Management (`env`)

The environments can be managed through the `env` commands.

 * `add` ([details](add)) - add a new environment
 * `info` ([details](info)) - show current environment information
 * `list` ([details](list)) - list all locally-configured environments
 * `remove` ([details](remove)) - remove an environment


## Options

 * **`url`** - a fully-qualified URI for where the server can be reached
 * `banner` - a custom message which may be shown to users
 * `name` - a default environment name which may be suggested to users
 * `title` a human friendly name for the environment which may be shown to users
 * `metadata` - a hash of arbitrary string keys and values which is opaque to the server and returned in the `/env/info` API endpoint
