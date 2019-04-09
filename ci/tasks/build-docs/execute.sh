#!/bin/sh

set -eu

task_dir="$PWD"

cd hugo-site

./bin/generate-metalink-artifacts-data.sh "file://$task_dir/artifacts/ssoca-final"

./bin/generate-repo-tags-data.sh "$task_dir/repo"
 # accidental lightweight tag; manual, pre-CI
echo '  v0.8.0: 2018-01-07 22:40:00 -0800' >> data/repo/tags.yml

mkdir -p static/img
wget -qO static/img/dpb587.jpg https://dpb587.me/images/dpb587-20140313a~256.jpg

latest_version=$( grep '^  ' data/repo/tags.yml | awk '{ print $1 }' | sed -e 's/^v//' -e 's/:$//' | sort -rV | head -n1 )

cat > config.local.yml <<EOF
title: ssoca
baseURL: "https://dpb587.github.io/ssoca"
googleAnalytics: "UA-37464314-3"
theme:
- balmy-release
- balmy
params:
  ThemeBrandIcon: /img/dpb587.jpg
  ThemeIncludeMenu: "/_menu"
  ThemeNavBadges: []
  ThemeNavItems:
  - title: docs
    url: /
  - title: releases
    url: /releases/
  - title: github
    url: "https://github.com/dpb587/ssoca"
  GitRepo: "https://github.com/dpb587/ssoca"
  GitEditPath: blob/master/docs
  GitCommitPath: commit
  releaseVersionLatest: "$latest_version"
EOF

hugo \
  --config="config.yml,config.local.yml" \
  --contentDir="$task_dir/repo/docs" \
  --destination="$task_dir/public"

./bin/git-commit.sh "$task_dir/public"
