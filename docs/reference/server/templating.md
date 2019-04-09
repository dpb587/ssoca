---
title: Templating
aliases:
- /server/templating/
---

# Templating

Some configuration values may support dynamic, runtime templating. Some services may use this to provide additional, user or request-specific customization. The following lookups are available in templated values.

Configuration settings which support templating are noted with `templatized` in documentation.


## Context

 * `Request` ([`net/http.Request`](https://golang.org/pkg/net/http/#Request)) -- the current request
 * `Token` -- the current authentication token
    * `Name` -- personal name (if available)
    * `Email` -- email address (if available)
    * `Username` -- username (if available)
    * `Groups` -- a list of scopes
       * `Contains(string)` -- check if a scope is available
       * `Matches(string)` -- check if a scope [matches](https://golang.org/pkg/path/filepath/#Match) the given pattern


## Examples

Using the [ssh service]({{< relref "../service/ssh" >}}) as an example:

```yaml
services:
- name: ssh
  type: ssh
  options:
    principals:
    # allow the mailbox of an email-based authentication
    - '{{ index ( split .Token.ID "@" ) 0 }}'
    # allow root if a member of a specific group
    - '{{ if .Token.Groups.Matches "bosh.4a12ea46-a526-4d32-8516-ed55feea5297.admin" }}jumpbox{{ end }}'
    # allow root if a member of any team in an adminorg
    - '{{ if .Token.Groups.Matches "adminorg/*" }}root{{ end }}'
    critical_options:
      # only allow connections from the client's current IP
      source-address: '{{ index ( split .Request.RemoteAddr ":" ) 0 }}/32'
    target:
      user: '{{ index ( split .Token.ID "@" ) 0 }}'
      host: 10.244.0.2
```
