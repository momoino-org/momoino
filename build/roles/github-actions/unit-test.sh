#!/bin/bash

set -eu

function main() {
  local cover_packages="x-operation/common/...,x-operation/console/...,x-operation/migration/..."
  local cover_profile="cover.profile"
  local testing_module="./testing/..."

  go run "github.com/onsi/ginkgo/v2/ginkgo@$GINKGO_VERSION" \
    -v \
    --procs=4 \
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
    --json-report=report.json \
    --timeout=5m \
    --poll-progress-after=120s \
    --poll-progress-interval=30s \
    --github-output \
    "$testing_module"
}

main "$@"
