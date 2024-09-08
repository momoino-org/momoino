#!/bin/bash

set -eu

function main() {
  local cover_packages="wano-island/common/...,wano-island/console/...,wano-island/migration/..."
  local cover_profile="coverage.txt"
  local testing_module="./..."

  go run github.com/onsi/ginkgo/v2/ginkgo \
    --randomize-all \
    --randomize-suites \
    --fail-on-pending \
    --fail-on-empty \
    --keep-going \
    --cover \
    --coverprofile="$cover_profile" \
    -coverpkg="$cover_packages" \
    --race \
    --trace \
    --timeout=5m \
    --poll-progress-after=120s \
    --poll-progress-interval=30s \
    --github-output \
    "$testing_module"
}

main "$@"
