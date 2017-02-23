#!/bin/bash

set -eu

cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )/../../.."

source .envrc

ginkgo -r -cover -covermode count

gover
GIT_BRANCH="$( git name-rev --name-only HEAD )" goveralls \
  -coverprofile=gover.coverprofile \
  -service concourse \
  -repotoken "$COVERALLS_TOKEN"
