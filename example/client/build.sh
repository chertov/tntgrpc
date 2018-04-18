#!/bin/bash
set -e
export CGO_ENABLED=0
export GOARCH=amd64
PATH=/usr/local/go/bin/:$PATH
export GOOS=linux
case "$OSTYPE" in
    solaris*) echo "SOLARIS" ;;
    darwin*)  export GOOS=darwin ;;
    linux*)   export GOOS=linux ;;
    bsd*)     echo "BSD" ;;
    msys*)    echo "WINDOWS" ;;
    *)        echo "unknown: $OSTYPE" ;;
esac
CURRENT_DIR=$(pwd)
GO_PROJECT_PATH=$CURRENT_DIR/$(dirname $0)
export GOPATH=$GO_PROJECT_PATH:$GO_PROJECT_PATH/app
MODULE_PATH=$GO_PROJECT_PATH/app/src/app
TARGET_PATH=$CURRENT_DIR/bin
PKG_PATH=$GO_PROJECT_PATH/pkg
echo GOPATH: $GOPATH
echo GO_PROJECT_PATH: $GO_PROJECT_PATH
echo MODULE_PATH: $MODULE_PATH
echo PKG_PATH: $PKG_PATH
echo TARGET_PATH: $TARGET_PATH
echo
cd $MODULE_PATH

mkdir -p $GO_PROJECT_PATH/app/src/helloworld
protoc --go_out=plugins=grpc:$GO_PROJECT_PATH/app/src/helloworld -I=/proto helloworld.proto
go get -t -d -v ./
echo

# Build for mac
export GOOS=darwin
TARGET=$TARGET_PATH/client_mac
echo Build for darwin...
echo GOOS: $GOOS
echo GOARCH: $GOARCH
echo TARGET: $TARGET
rm -f $TARGET
go build -o $TARGET -i -pkgdir $PKG_PATH

# Build for linux
export GOOS=linux
TARGET=$TARGET_PATH/client_linux
echo Build for linux...
echo GOOS: $GOOS
echo GOARCH: $GOARCH
echo TARGET: $TARGET
rm -f $TARGET
go build -o $TARGET -i -pkgdir $PKG_PATH

# Build for Windows 64
export GOOS=windows
TARGET=$TARGET_PATH/client.exe
echo Build for windows...
echo GOOS: $GOOS
echo GOARCH: $GOARCH
echo TARGET: $TARGET
rm -f $TARGET
go build -o $TARGET -i -pkgdir $PKG_PATH

cd $CURRENTDIR
