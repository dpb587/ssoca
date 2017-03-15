 * **Overview**
    * [Architecture](architecture)
    * [Installation](installation)
 * **[Authentication](auth/authn)** - how users authenticate
    * [Cloud Foundry UAA](auth/authn/uaa)
    * [GitHub](auth/authn/github)
    * [Google](auth/authn/google)
 * **[Authorization](auth/authz)** - how users are authorized
 * **[Certificate Authority](certauth)** - how certificates can be signed
    * [In-Memory](certauth/memory)
    * [Local Filesystem](certauth/fs)
 * **[Client](client)** - how end users interact with ssoca
 * **[Server](server)** - how the ssoca server runs
    * [Configuration](server/config)
    * [Frontend UI](server/ui)
    * [Logging](server/logging)
 * **Services** - how external services are used
    * [Authentication](service/auth)
    * [Document Root](service/docroot)
    * [Download](service/download)
    * [Environment](service/env)
    * [OpenVPN](service/openvpn)
    * [SSH](service/ssh)

There's also some [development notes](dev) if you're looking to make changes or work from source.
