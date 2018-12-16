#!/bin/bash

set -eu

cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )/../../.."

export GOPATH=$PWD/../../../..

exec ginkgo -r
