# Environment Management (`env`)

The environments can be managed through the `env` commands.

 * `info` ([details](info-cmd)) - show environment information
 * `list` ([details](list-cmd)) - list all locally-configured environments
 * `rename` ([details](rename-cmd)) - set a new name for the environment
 * `services` ([details](services-cmd)) - show current services available from the environment
 * `set` ([details](set-cmd)) - configure the connection to an environment
 * `set-option` ([details](set-option-cmd)) - set a local client option in the environment
 * `unset` ([details](unset-cmd)) - remove all configuration for an environment
 * `update-client` ([details](update-client-cmd)) - download the latest client from the environment


## Options

These options must be configured with the top-level `env` section rather than a service.

 * **`url`** - a fully-qualified URI for where the server can be reached
 * `banner` - a custom message which may be shown to users
 * `name` - a default environment name which may be suggested to users
 * `title` a human friendly name for the environment which may be shown to users
 * `metadata` - a hash of arbitrary string keys and values which is opaque to the server and returned in the `/env/info` API endpoint
