---
title: Download (download)
aliases:
- /service/download/
---

# Download (`download`)

A service to expose static files from the local filesystem which can be downloaded by a client.

*Tip*: A `download` service can be used in conjunction with the `env.update_service` setting which provides `ssoca` client binaries for users to download ([example](https://github.com/dpb587/ssoca-bosh-release/blob/46f9a6e0cc45cfbe0ed4ec4b14d155dbeee0c303/jobs/ssoca/templates/etc/server.conf.erb#L113-L125)).


## Client Commands

The published files can be accessed through `download` subcommands.

 * `get` ([details]({{< ref "/reference/service/download/get-cmd.md" >}})) -- get an artifact
 * `list` ([details]({{< ref "/reference/service/download/list-cmd.md" >}})) -- list available artifacts


## Server Configuration Options

 * **`glob`** -- a glob path for matching files to publish for downloads
