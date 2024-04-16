#!/bin/bash
CURRENT_DIR=$(pwd)

for module in $(find $CURRENT_DIR/protos/* -type d); do
    protoc -I /usr/local/include \
           -I $GOPATH/pkg/mod/github.com/gogo/protobuf@v1.3.2 \
           -I $CURRENT_DIR/protos/ \
            --gofast_out=plugins=grpc:$CURRENT_DIR/genproto/ \
            $module/*.proto;
done;

for module in $(find $CURRENT_DIR/genproto/* -type d); do
  if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i "" -e "s/,omitempty//g" $module/*.go
  else
    sed -i -e "s/,omitempty//g" $module/*.go
  fi
done;