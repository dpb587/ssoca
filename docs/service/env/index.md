# Environment Management (`env`)

The environments can be managed through the `env` commands.

 * `add` ([details](add-cmd)) - add a new environment
 * `info` ([details](info-cmd)) - show current environment information
 * `list` ([details](list-cmd)) - list all locally-configured environments
 * `remove` ([details](remove-cmd)) - remove an environment


## Options

These options must be configured with the top-level `env` section rather than a service.

 * **`url`** - a fully-qualified URI for where the server can be reached
 * `banner` - a custom message which may be shown to users
 * `name` - a default environment name which may be suggested to users
 * `title` a human friendly name for the environment which may be shown to users
 * `metadata` - a hash of arbitrary string keys and values which is opaque to the server and returned in the `/env/info` API endpoint
