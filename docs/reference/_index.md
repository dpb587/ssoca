---
title: Overview
type: docs
aliases:
- /architecture/
---

# Overview

Single Sign On for services that use Certificate Authority-based authentication.


## Architecture

There are three components which `ssoca` brings together, and each component can be configured with at least one of several backends.

1. **Certificate Authority / Provider** - responsible for signing certificates which `ssoca` has created/reviewed and approved. The certificates it generates will typically be very short-lived (e.g. expire within a minute) and restricted to the user's session (e.g. username, IP address). [Learn more](certauth).
1. **Identity Provider** - responsible for authenticating users and providing authorization details for their session. This may be handled through browser interaction or CLI prompts, but once completed the client will typically retain a medium-lived token (e.g. one day) with their authorization details. [Learn more](authn).
1. **Service Provider** - responsible for accepting certificates signed by a trusted CA as a form of authentication. Once authenticated, the original, authenticating certificate should no longer be needed. [Learn more](service).

The following diagram demonstrates the high-level interactions which occur between the components.

<div class="wsd" wsd_style="roundgreen"><pre>
  note over ssoca-cli,ssoca,cert-provider,identity-provider,service-provider
    Authentication
  end note

  ssoca-cli->ssoca: check authentication method
  ssoca->ssoca-cli: authentication method
  ssoca-cli->identity-provider: authenticate (via browser, token, ...)
  identity-provider->ssoca-cli: authentication token (24h)

  note over ssoca-cli,ssoca,cert-provider,identity-provider,service-provider
    Services
  end note

  ssoca-cli->ssoca: certificate signing request
  ssoca->ssoca: authn/authz verification
  ssoca->cert-provider: certificate signing request
  cert-provider->ssoca: certificate (2m)
  ssoca->ssoca-cli: certificate + service config
  ssoca-cli->service-provider: certificate-based connection
</pre></div>

<script type="text/javascript" src="https://www.websequencediagrams.com/service.js"></script>


## Related Articles

 * [Scalable and secure access with SSH](https://code.fb.com/security/scalable-and-secure-access-with-ssh/) - Facebook discussing some of the technical architecture behind a CA-based approach and how they use it for SSH authentication in their environments.
