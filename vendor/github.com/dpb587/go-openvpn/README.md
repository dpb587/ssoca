# go-openvpn

Some libraries for [OpenVPN](https://openvpn.net/) tools.


## Utilities


### `ovpn-to-onc`

Convert an OpenVPN profile `*.ovpn` file into an [ONC](https://chromium.googlesource.com/chromium/src/+/master/components/onc/docs/onc_spec.md) `*.onc` file for use in Chrome OS:

    $ ovpn-to-onc < vpn.ovpn > converted-vpn.onc


### `ovpn-to-json`

Convert a profile into JSON to split out directives and their arguments...

    $ ovpn-to-json < vpn.ovpn | jq -r .cipher[0]
    AES-256-CBC

    $ ovpn-to-json < vpn.ovpn | jq -r .key
    -----BEGIN RSA PRIVATE KEY-----
    MIIEowIBAAKCAQEA0haiWp2QxsJLsN2YkGiDUlT4CRxR95L8H6BkF/cla1uwZBJ9


## License

[MIT License](LICENSE)
