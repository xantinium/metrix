#!/bin/bash

(cd ./cmd/staticlint; go build .)
go vet -vettool=./cmd/staticlint/staticlint ./...
