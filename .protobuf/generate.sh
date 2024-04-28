#/bin/bash

## Generate multivac-edges service from remote url

mkdir -p "../services/multivac-edges"
wget "http://dev.multivac.studio/multivac-edges/service.proto" -O ./multivac-edges.proto
protoc --go_out=../services/multivac-edges --go_opt=paths=source_relative \
    --go-grpc_out=../services/multivac-edges --go-grpc_opt=paths=source_relative \
    ./multivac-edges.proto
rm ./multivac-edges.proto