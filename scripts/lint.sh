#!/bin/bash

CURRENT_DIR=$(cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd)
cd $(dirname $CURRENT_DIR)

(cd ./cmd/staticlint; go build .)
go vet -vettool=./cmd/staticlint/staticlint $(go list ./... | grep -v -E "/gen")
