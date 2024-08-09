#!/bin/bash

set -eu

# Main function to tidy up Go modules and synchronize the workspace.
function main() {
    # Tidy up Go modules in the common directory.
    cd ./common
    go mod tidy || { echo "go mod tidy failed in ./common"; exit 1; }
    echo "go mod tidy succeeded in ./common"

    # Tidy up Go modules in the console directory.
    cd ../console
    go mod tidy || { echo "go mod tidy failed in ../console"; exit 1; }
    echo "go mod tidy succeeded in ../console"

    # Tidy up Go modules in the migration directory.
    cd ../migration
    go mod tidy || { echo "go mod tidy failed in ../migration"; exit 1; }
    echo "go mod tidy succeeded in ../migration"

    # Tidy up Go modules in the mocks directory.
    cd ../testing
    go mod tidy || { echo "go mod tidy failed in ../testing"; exit 1; }
    echo "go mod tidy succeeded in ../testing"

    # Synchronize the Go workspace.
    cd ../
    rm -rf go.work.sum
    go work sync || { echo "go work sync failed"; exit 1; }
    echo "go work sync succeeded"
}

# Execute the main function with all arguments passed to the script.
main "$@"
