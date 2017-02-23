#!/bin/bash

set -eu

cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )/../../.."

source .envrc

exec ginkgo -r
