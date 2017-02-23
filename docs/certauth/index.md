# Certificate Authority

Services are generally configured to trust a specific certificate authority. These CAs may be backed by different providers and take responsibility for signing standard x.509 and OpenSSH certificates.

## Supported Backends

 * [In-Memory](memory)
 * [Local Filesystem](fs)


## Development

For development, you can easily create a self-signed certificate authority with [`certstrap`](https://github.com/square/certstrap)...

    $ certstrap init --common-name test-ca --passphrase ''
