#!/bin/bash

set -eu

function main() {
  kubectl port-forward --namespace automation-test-datastore svc/postgres-postgresql 5432:5432
}

main "$@"
