# Development

Some notes to remember when working from source...


## Commits

Before committing, run the `bin/pre-commit` script to...

 * regenerate fakes
 * `go fmt` source files
 * run all tests
 * build all binaries
 * regenerate client command docs

Review for unexpected changes before including them in your commit.


## Shortcuts

Some shortcuts instead of the built-in `go` commands...


### `go run`

 * `bin/client` - shortcut to run the client from source in any directory
 * `bin/server` - shortcut to run the server from source in any directory


### `go build`

Run `bin/build` to build both client and server for all supported architectures and operating systems. Optionally pass a version as the first argument (defaults to `0.0.0`). Binaries are put into the `tmp` directory of the repository and use the following naming convention:

    ssoca-(client|server)-$VERSION-$GOOS-$GOARCH
    # e.g. ssoca-client-0.1.0-darwin-amd64

    $ bin/pre-commit


## Certificate Authorities

You might find [certstrap](https://github.com/square/certstrap) useful for generating certificates.


### Installation

    $ go get github.com/square/certstrap
    $ pushd $GOPATH/src/github.com/square/cerstrap
    $ ./bin/build
    $ mv bin/certstrap $GOPATH/bin
    $ popd


### Usage

    $ certstrap --depot-path tmp/dev-server init --cn ca --passphrase ''
    $ chmod 0600 tmp/dev-server/ca.key
    $ ssh-keygen -f tmp/dev-server/ca.key -y > tmp/dev-server/ca-cert.pub
    $ certstrap --depot-path tmp/dev-server request-cert --cn server --domain localhost --passphrase ''
    $ chmod 0600 tmp/dev-server/server.key
    $ certstrap --depot-path tmp/dev-server sign server --CA ca


## Random Tips

 * pipe server log output through `jq` for more readable messages
