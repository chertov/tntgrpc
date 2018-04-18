#!/bin/bash
TNT_GRPC_NAME="grpc"
PROTOS=""
for i in "$@"; do
    key="$1"
    case $key in
        --name=*)
        TNT_GRPC_NAME="${i#*=}"
        shift # past --name=value
        continue
        ;;
        *.proto)
        PROTOS="$PROTOS /proto/$key"
        shift # past --name=value
        continue
        ;;
    esac
done

echo NAME: $TNT_GRPC_NAME
echo PROTOS: $PROTOS

export TNT_GRPC_NAME=$TNT_GRPC_NAME
/usr/local/bin/protoc --plugin=protoc-gen-grpc=/usr/local/bin/grpc_cpp_plugin --grpc_out=./gen/ -I=/proto/ --cpp_out=./gen/ $PROTOS \
&& bin/proto_tool $PROTOS \
&& cmake ./ \
&& make -j 6
