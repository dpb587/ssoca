---
title: Document Root (docroot)
aliases:
- /service/docroot/
---

# Document Root (`docroot`)

A service to expose files from the local filesystem (specifically, for exposing a simple web frontend).

This service does not provide any client-side commands. Instead, when users are viewing the service endpoint from a browser they will be able to view HTML files from the directory.

*Tip*: A `docroot` service can be used in conjunction with the `server.root_redirect` setting which provides a browser-friendly homepage for the server ([example](https://github.com/dpb587/ssoca-bosh-release/blob/46f9a6e0cc45cfbe0ed4ec4b14d155dbeee0c303/jobs/ssoca/templates/etc/server.conf.erb#L91-L111)).


## Server Configuration Options

 * **`path`** -- the directory which has browser frontend files and which will be served at the endpoint
