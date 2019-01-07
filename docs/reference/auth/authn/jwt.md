# JSON Web Tokens

Some authentication providers will use [JWT](https://jwt.io/) to create a signed representation of a user's authorization details. These tokens are typically valid for a relatively short period (24 hours) before the user needs to re-authenticate with the identity provider for an updated token.


## Options

 * **`private_key`** -- a PEM-formatted private key
 * `validity` -- a [duration](https://golang.org/pkg/time/#ParseDuration) for how long authentication tokens will be remembered (default `24h`)


## General Notes

You can use the following to generate a new key for testing...

    openssl genrsa -out jwt.key 2048
