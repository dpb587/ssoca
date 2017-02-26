# Logging

The server process emits logs to standard output in JSON format.


## Global Fields

The following fields will be present on every message.

 * `level` - severity of the message (i.e. `DEBUG`, `INFO`, `WARN`, `ERROR`, `FATAL`, `PANIC`)
 * `msg` - a brief summary describing the event which occurred (e.g. `Signed ssh certificate`)
 * `time` - time the message was logged in ISO8601 format (e.g. `2017-02-21T02:12:22Z`)


## HTTP Fields

The following fields will be present on messages related to an HTTP request or response.

 * `server.request.id` - a UUID identifying the request (this may appear in other messages for correlation when they're initiated by an HTTP request)
 * `server.request.method` - the HTTP method used (e.g. `POST`)
 * `server.request.path` - the URL path requested (e.g. `/ssh/sign-public-key`)
 * `server.request.remote_addr` - the remote address of the connecting client (e.g. `[::1]:50520`)


## Services

The following fields will be present on messages related to a specific service.

 * `service.name` - the configured service name (e.g. `failover-vpn`)
 * `service.type` - the service type (e.g. `openvpn`)


## Authentication

The following fields will be present on messages which were being performed by an authenticated user.

 * `auth.id` - the authenticated user (e.g. `somebody@example.com`)


## Certificate Authority

The following fields will be present on messages related to certificate operations.

 * `certauth.name` - the configured certificate authority name (e.g. `vpn`)
 * `certauth.ssh.key_id` - the identifying key ID of an SSH certificate
 * `certauth.ssh.valid_after`, `certauth.ssh.valid_before` - the validity range of an SSH certificate in ISO8601 format
 * `certauth.x509.serial` - the identifying serial number of an x509 certificate
 * `certauth.x509.not_before`, `certauth.x509.not_after` - the validity range of an x509 certificate in ISO8601 format
 * `certauth.x509.common_name` - the identifying common name of an x509 certificate
