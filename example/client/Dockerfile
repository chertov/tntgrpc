FROM alpine:3.5
MAINTAINER chertovmv@gmail.com

RUN set -x \
    && apk add --no-cache --virtual .build-deps \
        git curl gcc g++ autoconf automake libtool make go

RUN set -x \
    && mkdir -p /usr/src/grpc \
    && git clone --progress --verbose -b $(curl -L https://grpc.io/release) https://github.com/grpc/grpc /usr/src/grpc \
    && git -C /usr/src/grpc submodule update --init \
    && make -C /usr/src/grpc \
    && make -C /usr/src/grpc install \
    && make -C /usr/src/grpc/third_party/protobuf install \
    && rm -rf /usr/src/grpc

RUN export GOPATH=/usr/local/go && go get -u github.com/golang/protobuf/proto \
    && go get -u github.com/golang/protobuf/protoc-gen-go \
    && go get -u google.golang.org/grpc

WORKDIR /client_src
# COPY ./helloworld.proto ./
# COPY ./test.proto ./
# COPY ./client/app/src/app/client.go ./app/src/app/client.go
# COPY ./client/build.sh ./
# RUN ls ./
# RUN mkdir -p /output
# RUN ./build.sh

# CMD ["mv", "/output", "/hostoutput/bin"]
CMD ["./build.sh"]
# CMD ["rm -rf /hostoutput/bin ; mv /output /hostoutput/bin"]
