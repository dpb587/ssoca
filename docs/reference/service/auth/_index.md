---
title: Authentication (auth)
aliases:
- /service/auth/
---

# Authentication (`auth`)

Several different types of authentication are supported and the server will configure whichever is appropriate for the environment.

 * [Cloud Foundry UAA]({{< relref "../uaa-auth" >}})
 * [GitHub]({{< relref "../github-auth" >}})
 * [Google]({{< relref "../google-auth" >}})

## Client Commands

Authentication workflows can be managed through `auth` subcommands.

 * `info` ([details]({{< ref "/reference/service/auth/info-cmd.md" >}})) -- show current authentication information
 * `login` ([details]({{< ref "/reference/service/auth/login-cmd.md" >}})) -- authenticate for a new token
 * `logout` ([details]({{< ref "/reference/service/auth/logout-cmd.md" >}})) -- revoke an authentication token


## Client Options

For advanced customization, the following options may be configured for an environment to influence how authentication workflows operate.


### `bind` (Web-based Bind Address)

If the authentication service needs to start a local web server during authentication, by default, a random port will be bound on `localhost`. This behavior can be overridden with the `bind` option to specify a specific IP or port.

For example, to force binding to port `8085` to enable static tunneling configuration, you might use...

```bash
ssoca env set-option auth.bind "localhost:8085"
```

### `open_command` (Interactive Login)

If the user needs to visit a URL during authentication, the CLI will attempt to open the URL automatically. By default, the system's open command is invoked, but this can be overridden with the `open_command` option if advanced usage is required. The URL will be appended to the command.

For example, to open the URL in Google Chrome with a specific profile for the user, you might use...

```bash
ssoca env set-option auth.open_command "[ sudo, -u, $USER, /Applications/Google Chrome.app/Contents/MacOS/Google Chrome, --profile-directory=Default, --disable-gpu ]"
```

{{< note type="info" >}}
  Shell environment variables are not interpolated by the client at runtime (i.e. `$USER` is parsed by shell here). The `sudo` usage allows the process to interact with the user's console UI if the command is run as `root`. The `--disable-gpu` suppresses a seemingly innocuous warning message.
{{< /note >}}
