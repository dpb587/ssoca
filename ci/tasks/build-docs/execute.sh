#!/bin/sh

set -eu

mkdir -p hugo-site/data

meta4-repo filter --format=json file://$PWD/artifacts/ssoca-final > hugo-site/data/releaseArtifacts.json

(
  cd repo
  git fetch --tags $( git remote get-url origin | sed 's#git@github.com:#https://github.com/#' ) # avoid stale concourse resource caches

  echo 'dates:'
  echo '  v0.8.0: 2018-01-07 22:40:00 -0800' # lightweight tag; manual, pre-CI
  git log --tags --simplify-by-decoration --pretty="format:%D: %ai" | sed 's#^HEAD -> master, ##' | sed 's#, origin/master:#:#' | grep -E '^tag: [^ ]+:' | sed 's#^tag: #  #'
) > hugo-site/data/repositoryTags.yml

cd hugo-site

hugo --contentDir=../repo/docs

cd ..

mv hugo-site/public/* public/

cd public

git config --global user.email "${git_user_email:-ci@localhost}"
git config --global user.name "${git_user_name:-CI Bot}"
git init
git add .
git commit -m 'build docs'
