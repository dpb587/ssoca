#!/bin/bash

set -eu

task_dir="$PWD"
release_dir="$task_dir/release"
version="$( cat version/version )"

cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )/../../.."

build_dir="$task_dir/build"
metalink_path="$build_dir/v$version.meta4"

mkdir "$release_dir/blobs"

if [ -e "$task_dir/repo/docs/release-notes/v${version}.md" ]; then
  sed -E 's/^# .+//' "$task_dir/repo/docs/release-notes/v${version}.md" > $release_dir/notes.md
  echo "" >> $release_dir/notes.md
fi

(
  echo "Artifacts"
  echo ""
  echo "                                                              sha256 file"
) >> $release_dir/notes.md

for file in $( meta4 files --metalink=$metalink_path ); do
  cp "$build_dir/$file" "$release_dir/blobs/$file"

  echo "    $( meta4 file-hash --metalink=$metalink_path --file=$file sha-256 )" "$file" >> $release_dir/notes.md
done

( cd "$task_dir/repo" ; git rev-parse HEAD ) > "$release_dir/commit"
echo "$version" > $release_dir/version
echo "v$version" > $release_dir/tag
echo "Release v$version" > $release_dir/title
