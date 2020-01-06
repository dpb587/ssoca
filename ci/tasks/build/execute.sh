#!/bin/bash

set -eu

build_dir="$PWD/build"
version="$( cat version/version )"

cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )/../../.."

export GOPATH=$PWD/../../../..
export PATH="$GOPATH/bin:$PATH"

./bin/install-tools
./bin/build "$version"

mv tmp/ssoca-client-* tmp/ssoca-server-* "$build_dir/"
