#!/bin/bash

set -eu

build_dir="$PWD/build"
version="$( cat version/version )"

cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )/../../.."

source .envrc

./bin/build "$version"

mv tmp/ssoca-client-* tmp/ssoca-server-* "$build_dir/"
