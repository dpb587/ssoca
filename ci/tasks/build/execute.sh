#!/bin/bash

set -eu

build_dir="$PWD/build"
version="$( cat version/version )"

s3_prefix="${s3_prefix:-}"

export AWS_ACCESS_KEY_ID="$s3_access_key"
export AWS_SECRET_ACCESS_KEY="$s3_secret_key"

cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )/../../.."

source .envrc

./bin/build "$version"

mv tmp/ssoca-client-* tmp/ssoca-server-* "$build_dir/"
