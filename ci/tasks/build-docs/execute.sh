#!/bin/sh

set -eu

task_dir="$PWD"

cd hugo-site

./bin/generate-metalink-artifacts-data.sh "file://$task_dir/artifacts/ssoca-final"

./bin/git-fetch.sh "$task_dir/repo" # avoid stale concourse resource caches
./bin/generate-repo-tags-data.sh "$task_dir/repo"
 # accidental lightweight tag; manual, pre-CI
echo '  v0.8.0: 2018-01-07 22:40:00 -0800' >> data/repo/tags.yml

latest_version=$( grep '^  ' data/repo/tags.yml | awk '{ print $1 }' | sed -e 's/^v//' -e 's/:$//' | sort -rV | head -n1 )

cat > config.local.yml <<EOF
baseURL: "https://dpb587.github.io/ssoca"
googleAnalytics: "UA-37464314-3"
params:
  releaseVersionLatest: "$latest_version"
EOF

hugo \
  --config="config.yml,$task_dir/repo/docs/config.yml,config.local.yml" \
  --contentDir="$task_dir/repo/docs" \
  --destination="$task_dir/public"

./bin/git-commit.sh "$task_dir/public"
