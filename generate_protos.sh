#!/bin/bash
# chmod +x generate_protos.sh

set -e

PROTO_DIR="./"

OUT_DIR="./genprotos"

mkdir -p $OUT_DIR

protoc -I=$PROTO_DIR \
  --go_out=$OUT_DIR \
  --go-grpc_out=$OUT_DIR \
  $PROTO_DIR/judge_service_protos/judge_service.proto

echo "protos generated successfully"
