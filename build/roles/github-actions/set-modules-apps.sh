#/bin/bash

set -eu

function main() {
  local modules=$(go list -m -json | jq -s '.' | jq -c '[.[].Dir]')
  local count=$(echo "$modules" | jq length)
  local -a apps=()

  for ((i = 0; i < count; i++)); do
    path=$(echo $modules | jq -r ".[$i]")

    if [ -f "$path/Dockerfile" ]; then
      apps+=("$path")
    fi
  done

  echo "apps=$(jq -c -n '$ARGS.positional' --args "${apps[@]}")" >>$GITHUB_OUTPUT
  echo "modules=$modules" >>$GITHUB_OUTPUT
}

main "$@"
