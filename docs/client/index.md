# Client

Most users will be executing the `ssoca` client binary from their workstations.


## First-time Setup

Users should be given a URL for where to find the `ssoca` server. If they open the URL from a browser, they can download

    $ wget -O- https://ssoca.example.com/download/useragent-bundle | tar -xzf-
    $ cd ssoca.example.com
    $ direnv allow

Then `auth login` can be used to to login interactively using whatever authentication method has been configured by the server (typically a URL will be shown)...

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


Or directly work with a service (some services may provide default endpoints making arguments optional)...

    $ ssoca http curl https://app.acme-dev.example.com/
    $ ssoca ssh -s jumpbox connect
    $ ssoca openvpn -s vpn connect
    $ ssoca sshuttle connect
