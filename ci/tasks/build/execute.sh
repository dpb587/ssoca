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

metalink_path="$build_dir/v$version.meta4"

meta4 create --metalink="$metalink_path"
meta4 set-published --metalink="$metalink_path" "$( date -u +%Y-%m-%dT%H:%M:%SZ )"

cd tmp

for file in $( find . -type f -perm +111 -maxdepth 1 -name "ssoca-*-$version-*" | cut -c3- | sort ); do
  echo "$file"

  meta4 import-file --metalink="$metalink_path" --file="$file" --version="$version" "$file"

  if [ -n "$s3_host" ]; then
    sha1=$( meta4 file-hash --metalink=$metalink_path --file="$file" sha-1 )
    meta4 file-upload --metalink="$metalink_path" --file="$file" "$file" "s3://$s3_host/$s3_bucket/${s3_prefix}v$version/$sha1"
  fi

  mv "$file" "$build_dir/$file"
done

cat "$metalink_path"
