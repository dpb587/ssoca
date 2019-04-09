---
title: Client Installation
aliases:
- /client/
---

# Client Installation

Most users will be executing the `ssoca` client binary from their workstations which require an initial setup.


## Dependencies

First, users must install the `ssoca` client binary. After you install using one of the following methods, run `ssoca version` to verify. This only needs to happen once per workstation.


### Homebrew or Linuxbrew

For users of [Homebrew](https://brew.sh/) (macOS) or [Linuxbrew](http://linuxbrew.sh/), you may use the [dpb587/homebrew-tap](https://github.com/dpb587/homebrew-tap) tap to install the latest official binaries.

```bash
brew install dpb587/tap/ssoca
```


### Official Releases

Official binaries are listed on the [Official Releases]({{% ref "/releases/_index.md" %}}) page (you may also find the same artifacts from [GitHub Releases](https://github.com/dpb587/ssoca/releases)). Find the correct `ssoca-client-*` for your operating system and platform and install it.

For example, with the latest release on macOS:

{{< release/download-install-artifact file="ssoca-client-.+-darwin-amd64" install="/usr/local/bin/ssoca" >}}


### Local Environment

Your local environment may also provide binaries for you to download as well. Visit your ssoca server from a browser for download links and checksums.


## Environment Configuration

Once you have the `ssoca` client available, you should configure your environment with an alias. Visiting the ssoca server from a browser may provide you with similar setup instructions. If your environment is using a custom CA certificate, use the `--ca-cert` option. This only needs to happen once per environment.

```bash
ssoca -e example-prod env set https://prod.example.com
```

You will receive a confirmation once it has connected successfully, then you may authenticate to verify access.

```bash
ssoca -e example-prod auth login
```

After authenticating, you can use one of the services provided by the server (e.g. OpenVPN or SSH).
