---
title: Environment (env)
aliases:
- /service/env/
---

# Environment (`env`)

Environments are used to distinguish between multiple *ssoca* installations. Aliases are used to reference environments in a more memorable way and are specified through the `-e`/`--environment` option or the `SSOCA_ENVIRONMENT` environment variable.


## Client Commands

The environments can be managed through `env` subcommands.

 * `info` ([details]({{< ref "/reference/service/env/info-cmd" >}})) -- show environment information
 * `list` ([details]({{< ref "/reference/service/env/list-cmd" >}})) -- list all locally-configured environments
 * `rename` ([details]({{< ref "/reference/service/env/rename-cmd" >}})) -- set a new name for the environment
 * `services` ([details]({{< ref "/reference/service/env/services-cmd" >}})) -- show current services available from the environment
 * `set` ([details]({{< ref "/reference/service/env/set-cmd" >}})) -- configure the connection to an environment
 * `set-option` ([details]({{< ref "/reference/service/env/set-option-cmd" >}})) -- set a local client option in the environment
 * `unset` ([details]({{< ref "/reference/service/env/unset-cmd" >}})) -- remove all configuration for an environment
 * `update-client` ([details]({{< ref "/reference/service/env/update-client-cmd" >}})) -- download the latest client from the environment


## Server Configuration Options

These options must be configured with the top-level `env` section rather than a service.

 * **`url`** -- a fully-qualified URI for where the server can be reached
 * `banner` -- a custom message which may be shown to users when they first connect
 * `name` -- a default alias name which may be suggested to users
 * `title` a human friendly name for the environment which may be shown to users
 * `metadata` -- a hash of arbitrary string keys and values which is opaque to the server and returned in the `/env/info` API endpoint
 * `update_service` -- the name of a configured `download` service which can provide `ssoca` client binaries
