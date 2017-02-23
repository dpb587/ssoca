Use [certstrap](https://github.com/square/certstrap) to generate certificates...

    $ go get github.com/square/certstrap
    $ pushd $GOPATH/src/github.com/square/cerstrap
    $ ./bin/build
    $ mv bin/certstrap $GOPATH/bin
    $ popd

Generate some certificates...

    $ certstrap --depot-path tmp/server init --cn ca --passphrase ''
    $ chmod 0600 tmp/server/ca.key
    $ ssh-keygen -f tmp/server/ca.key -y > tmp/server/ca-cert.pub
    $ certstrap --depot-path tmp/server request-cert --cn server --domain localhost --passphrase ''
    $ chmod 0600 tmp/server/server.key
    $ certstrap --depot-path tmp/server sign server --CA ca
    $ cat > tmp/server/config.yml <<EOF
    auth:
      type: http
      options:
        users:
          - username: admin
            password: nimda
    certauths:
      - name: default
        type: fs
        options:
          private_key_path: tmp/server/ca.key
    server:
      certificate_path: tmp/server/server.crt
      private_key_path: tmp/server/server.key
    services:
      - name: ssh
        type: ssh
        options:
          certauth: default
    EOF

Start the server...

    $ go run cli/server/main.go

Run the client...

    $ go run cli/client/main.go -e localhost env add --ca-cert tmp/server/ca.crt https://localhost:18705

Run the tests...

    $ ginkgo -r

Lint before committing...

    $ bin/pre-commit

# Snippets

Convert PEM to OpenSSH public key...

    ssh-keygen -f ca-private.pem -y > ca-cert.pub

Investigate a signed public key...

    ssh-keygen -L -f id_rsa-cert.pub

JWT key...

    openssl genrsa -out jwt.key 2049

More...

    ssh-keygen -C CA -f ca
    ssh-keygen -t ecdsa -f jdoe

    cat ~/.ssh/id_rsa.pub | jq -sR '{"public_keys":[.]}' | curl -kd@- https://127.0.0.1:18705/ssh/sign
