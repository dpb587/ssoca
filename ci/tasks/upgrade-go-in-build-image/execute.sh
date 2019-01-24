#!/bin/sh

set -eu

git clone file://$PWD/repo repo-output

dockerfile=ci/images/build/Dockerfile
version=$( cat golang/.resource/version )

cd repo-output

sed -ri "s/^FROM golang:.+$/FROM golang:$version-stretch/" "$dockerfile"

git add "$dockerfile"

if git diff --staged --exit-code --quiet ; then
  # no changes pending
  exit
fi

git config --global user.email "${git_user_email:-ci@localhost}"
git config --global user.name "${git_user_name:-CI Bot}"

git commit -m "ci: upgrade build image to golang/$version"
