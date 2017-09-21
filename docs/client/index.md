# Client

Most users will be executing the `ssoca` client binary from their workstations.


## First-time Setup

Users can download the `ssoca` client binary from the project's [latest release](https://github.com/dpb587/ssoca/releases/latest)...

    $ wget -O /usr/local/bin/ssoca https://github.com/dpb587/ssoca/releases/download/v0.7.0/ssoca-client-0.7.0-darwin-amd64
    $ echo "34f8334120adc3028b685703531abc4044cb454f815c53f0f0a8cf85e86c07fb  /usr/local/bin/ssoca" | shasum -a 256 -c
    $ chmod +x /usr/local/bin/ssoca

Then `auth login` can be used to login interactively using whichever authentication method has been configured by the server (typically a URL will be shown)...

    $ ssoca auth login
    Visit the following link to receive an authentication token...

      https://ssoca.acme-dev.example.com/auth/initiate

    token: ...paste...

Once authenticated, the user can review the available services...

    $ ssoca env info
    Service   Type      Metadata  
    auth      uaa_auth  -
    jumpbox   ssh       -
    sshuttle  sshuttle  -
    vpn       openvpn   -

And then directly use a service. For example...

    $ ssoca openvpn -s vpn exec --sudo
