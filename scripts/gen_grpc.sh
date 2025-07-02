#!/bin/bash

CURRENT_DIR=$(cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd)
cd $(dirname $CURRENT_DIR)

RPC_PROTO_PATH=internal/presentation/rpc/proto/*.proto
RPC_PROTO_MODULE=github.com/xantinium/metrix/internal/presentation/rpc/gen
RPC_GEN_OUTPUT_PATH=./internal/presentation/rpc/gen

echo "Generating gRPC..."
GEN_RPC_ERR=$(protoc \
    --go_out=$RPC_GEN_OUTPUT_PATH \
    --go_opt=module=$RPC_PROTO_MODULE \
    --go-grpc_out=$RPC_GEN_OUTPUT_PATH \
    --go-grpc_opt=module=$RPC_PROTO_MODULE \
    $RPC_PROTO_PATH  2>&1 > /dev/null)

if [ $? != 0 ]; then
    echo "Failed to generate:"
    echo $GEN_RPC_ERR
    exit 1
fi
