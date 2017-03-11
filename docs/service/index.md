# Services

Services provide the client+server bindings between `ssoca` and the certificate, identity, or user service providers. The following services provide the user with certificate-based authentication to external services:

 * [HTTP](http)
 * [OpenVPN](openvpn)
 * [SSH](ssh)

Several meta-services provide commands and endpoints for authentication and additional resources:

 * [Authentication](auth)
 * [Document Root](docroot)
 * [Download](download)
 * [Environment](env)


## Configuration

Services require additional server-side configuration ([learn more](../server)). Configuration for services are documented on their respective service page.


### Templated Values

Some configuration values may support dynamic, runtime templating. Some services may use this to provide additional, user or request-specific customization. The following lookups are available in templated values.

 * `.Request.ClientIP` - the remote IP address of the client making the API request
 * `.Service.Name` - the name of the service
 * `.Token.Authenticated` - a boolean value indicating that the user is authenticated
 * `.Token.Username` - the username of the authenticated user
 * `.Token.Scopes` - a string array of authorized scopes

Configuration settings which support templating are noted with `templated` in documentation.
