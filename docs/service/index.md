# Services

Services can provide internal and external authentication behaviors.

 * [Document Root](service/docroot)
 * [Download](service/download)
 * [HTTP](service/http)
 * [OpenVPN](service/openvpn)
 * [SSH](service/ssh)


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
