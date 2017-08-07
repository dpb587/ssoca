# ssoca

[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Alpha stability](https://img.shields.io/badge/stability-alpha-yellow.svg)](#development)
[![Coverage Status](https://coveralls.io/repos/github/dpb587/ssoca/badge.svg?branch=master)](https://coveralls.io/github/dpb587/ssoca)

SSO for services that use CA-based authentication.

For when you might want...

 * ssh users to authenticate against Cloud Foundry UAA,
 * Google Cloud project owners to have access to an OpenVPN server, or
 * a GitHub team to access a network with sshuttle

With the caveat that this repo...

 * is a work in progress, and
 * is open source to help facilitate demos, discussion, and reviews to continue its evolution


## Summary

Supporting services like...

 * HTTP x.509 ([rfc5280](https://tools.ietf.org/html/rfc5280)) (in theory)
 * [OpenSSH](https://www.openssh.com/) ([rfc6187](https://tools.ietf.org/html/rfc6187))
 * [OpenVPN](https://openvpn.net/)
 * [sshuttle](https://github.com/apenwarr/sshuttle)

Supporting authentication from (and restricting by)...

 * [GitHub](https://github.com/) - organization, team, user
 * [Google](https://www.google.com/) - email, email domain, Cloud project+role
 * HTTP Basic

Supporting certificate authority keys stored in...

 * In-memory
 * Local filesystem

Supported technically by...

 * authentication being delegated to an external service (like Okta, UAA, GitHub, OAuth), and
 * external services being configured to trust a particular certificate authority, with
 * `ssoca` validating authentication and signing short-lived certificates.


## Details

 * [User Documentation](docs)
 * [Technical Documentation](https://godoc.org/github.com/dpb587/ssoca)
 * [BOSH Release](https://github.com/dpb587/ssoca-bosh-release)
 * [Roadmap](https://trello.com/b/LEu5Crqw/ssoca)
 * ssoca (s&#x014D;s&#x0259;, SO-sa)


## License

[MIT License](LICENSE)
