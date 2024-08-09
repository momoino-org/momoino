#!/bin/bash

set -eu

function main() {
  local modules="$1"
  local count=$(echo $modules | jq length)

  for ((i = 0; i < count; i++)); do
    local module_path=$(echo $modules | jq -r ".[$i]")

    cd "$module_path"

    go mod tidy || {
      echo "go mod tidy failed in $module_path"
      exit 1
    }

    git diff --exit-code --quiet go.mod go.sum || {
      echo "git diff failed in $module_path"
      exit 1
    }

    echo "go mod tidy succeeded in $module_path"
  done
}

main "$@"
