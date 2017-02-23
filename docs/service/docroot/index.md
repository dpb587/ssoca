# Document Root (`docroot`)

A service to expose files from the local filesystem (specifically, for exposing a simple web frontend).

This service does not provide any client-side commands. Instead, when users are viewing the service endpoint from a browser they will be able to view HTML files from the directory. A `docroot` service can be used in conjunction with the `server.root_redirect` setting which provides a browser-friendly homepage for the server.


## Options

 * **`path`** - the directory which has browser frontend files and which will be served at the endpoint

    path: /some/path/to/docroot
