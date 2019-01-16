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

    brew install dpb587/tap/ssoca


### GitHub Release

Official binaries and their checksums may also be found from the [Official Releases]({{% ref "/release" %}}). Find the correct `ssoca-client-*` for your operating system and platform.

{{% install/download-binary-steps %}}


### Local Environment

Your local environment may also provide binaries for you to download as well. Visit your ssoca server from a browser for download links and checksums.


## Environment Configuration

Once you have the `ssoca` client available, you should configure your environment with an alias. Visiting the ssoca server from a browser may provide you with similar setup instructions. If your environment is using a custom CA certificate, use the `--ca-cert` option. This only needs to happen once per environment.

    ssoca -e example-prod env set https://prod.example.com

You will receive a confirmation once it has connected successfully, then you may authenticate to verify access.

    ssoca -e example-prod auth login

After authenticating, you can use one of the services provided by the server (e.g. OpenVPN or SSH).
