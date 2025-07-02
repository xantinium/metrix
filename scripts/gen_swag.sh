#!/bin/bash

CURRENT_DIR=$(cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd)
cd $(dirname $CURRENT_DIR)

swag init -g ./internal/presentation/rest/handlers/handlers.go --ot "yaml" -o ./swagger
