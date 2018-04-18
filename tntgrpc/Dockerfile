FROM alpine:3.5
MAINTAINER chertovmv@gmail.com

ENV TARANTOOL_VERSION=2.0.4
ENV TARANTOOL_DOWNLOAD_URL=https://github.com/tarantool/tarantool.git

RUN set -x \
    && apk add --no-cache --virtual .build-deps \
        git curl gcc g++ autoconf automake libtool make cmake \
        nodejs go lua-dev

RUN set -x \
    && mkdir -p /usr/src/grpc \
    && git clone --progress --verbose -b $(curl -L https://grpc.io/release) https://github.com/grpc/grpc /usr/src/grpc \
    && git -C /usr/src/grpc submodule update --init \
    && make -C /usr/src/grpc \
    && make -C /usr/src/grpc install \
    && make -C /usr/src/grpc/third_party/protobuf install \
    && rm -rf /usr/src/grpc

RUN set -x \
    && apk add --no-cache --virtual .tarantool-build-deps \
        readline-dev libressl-dev yaml-dev curl-dev icu-dev \
    && mkdir -p /usr/src/tarantool \
    && git clone --progress --verbose "$TARANTOOL_DOWNLOAD_URL" /usr/src/tarantool \
    && git -C /usr/src/tarantool checkout "$TARANTOOL_VERSION" \
    && git -C /usr/src/tarantool submodule update --init --recursive \
    && (cd /usr/src/tarantool; \
        cmake   -DCMAKE_BUILD_TYPE=RelWithDebInfo \
                -DENABLE_BUNDLED_LIBYAML:BOOL=OFF \
                -DENABLE_BACKTRACE:BOOL=OFF \
                -DENABLE_DIST:BOOL=OFF \
                .) \
    && make -C /usr/src/tarantool \
    && make -C /usr/src/tarantool install \
    && make -C /usr/src/tarantool clean \
    && rm -rf /usr/src/tarantool \
    && : "---------- remove build deps ----------" \
    && apk del .tarantool-build-deps

WORKDIR /tntgrpc

RUN set -x \
    && mkdir -p /output/ \
    && mkdir -p ./gen

COPY  ./ ./
RUN (cd jstool && npm i && node generrors.js) && proto_tool/build.sh

# Example:
# RUN ./run.sh --name=myprotolib helloworld.proto test.proto
ENTRYPOINT ["./run.sh"]
CMD []
