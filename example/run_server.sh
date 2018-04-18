rm -rf mytntgrpclib.so; docker rm mytarantool

docker run -it -v $(pwd):/output -v $(pwd):/proto tntgrpc \
    --name=mytntgrpclib helloworld.proto

docker run --name mytarantool -p3301:3301 -p50051:50051 -it \
    -v $(pwd):/opt/tarantool \
    tarantool/tarantool:2 \
    tarantool /opt/tarantool/app.lua
