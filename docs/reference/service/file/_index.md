---
title: File (file)
aliases:
- /service/download/
- /reference/service/download/
---

# File (`file`)

A service to expose static files from the local filesystem which can be accessed by a client.

{{< note type="success" >}}
  A `file` service can be used in conjunction with the `env.update_service` setting which provides `ssoca` client binaries for users to download ([example](https://github.com/dpb587/ssoca-bosh-release/blob/46f9a6e0cc45cfbe0ed4ec4b14d155dbeee0c303/jobs/ssoca/templates/etc/server.conf.erb#L113-L125)).
{{< /note >}}


## Client Commands

The published files can be accessed through `file` subcommands.

 * `exec` ([details]({{< ref "/reference/service/file/exec-cmd.md" >}})) -- temporarily get and then execute a file
 * `get` ([details]({{< ref "/reference/service/file/get-cmd.md" >}})) -- download a file and verify its checksum
 * `list` ([details]({{< ref "/reference/service/file/list-cmd.md" >}})) -- list available files


## Server Configuration Options

 * **`glob`** -- a glob path for matching files to publish for downloads
