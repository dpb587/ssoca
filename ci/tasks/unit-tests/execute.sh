#!/bin/bash

set -eu

cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )/../../.."

export GOPATH=$PWD/../../../..
export PATH="$GOPATH/bin:$PATH"

./bin/install-tools

exec ginkgo -r
