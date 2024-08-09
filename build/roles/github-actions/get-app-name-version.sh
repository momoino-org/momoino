#!/bin/bash

set -eu

function main() {
  local app_path="$1"

  echo "app_name=$(basename $app_path)" >>$GITHUB_OUTPUT
  echo "app_version=$(cat $app_path/VERSION)" >>$GITHUB_OUTPUT
}

main "$@"
