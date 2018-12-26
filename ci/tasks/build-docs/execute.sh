#!/bin/sh

set -eu

mkdir -p hugo-site/repo hugo-site/data
mv repo/docs hugo-site/repo/docs
meta4-repo filter --format=json file://$PWD/artifacts/ssoca-final > hugo-site/data/releaseArtifacts.json

cd hugo-site
hugo

cd ..

mv hugo-site/public/* public/

cd public

git config --global user.email "${git_user_email:-ci@localhost}"
git config --global user.name "${git_user_name:-CI Bot}"
git init
git add .
git commit -m 'build docs'
