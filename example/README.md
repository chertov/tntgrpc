This example explains how to use GRPC generator for Tarantool application server.

## Start tarantool GRPC server
Remove the grpc library 'mytntgrpclib.so' and the docker container if exists.
```sh
rm -rf mytntgrpclib.so; docker rm mytarantool
```

On this step we are trying to generate mytntgrpclib.so from helloworld.proto and test.proto files.
We need only *.proto files to generate native library 'mytntgrpclib.so'.
We can to use docker container from 'tntgrpc' image for this.
We can specify the library name ('mytntgrpclib' in this example).
Also we must share two directories with proto files and output path for the container.
'/proto' path is a directory where grpc tool take proto files (helloworld.proto)
'/output' path is a directory where grpc tool place tarantool native library (mytntgrpclib.so)
```sh
docker run -it \
    -v $(pwd):/proto \  # our proto files is here
    -v $(pwd):/output \ # output files
    tntgrpc \
    --name=mytntgrpclib helloworld.proto
```
Ok, Now we have 'mytntgrpclib.so' and we can use it for GRPC service inside Tarantool.
We need to share port 3301 for tarantool and 50051 for grpc,
also we must share current folder for tarantool because lua code (app.lua) and grpc lib (mytntgrpclib.so) are here.
on this step we don't need *.proto files anymore. Only lua code and native library are required.
If you need you can use it with different way of course.
```sh
docker run --name mytarantool -it \
    -p3301:3301 -p50051:50051 \
    -v $(pwd):/opt/tarantool \
    tarantool/tarantool:2 \
    tarantool /opt/tarantool/app.lua
```

## Start GRPC client

Ok, we should build GRPC client for Tarantool GRPC service.
Let's build it with docker. It's supereasy!

We have Dockerfile with build tools, grpc and golang compilers. Let's build it.
```sh
docker build -t tnt_grpc_client_builder \
    -f ./client/Dockerfile ./
```
Let's start 'tnt_grpc_client_builder' container.
We need to share *.proto files and Goland source code for it.
We will have been the clients for Mac, Linux and Windows in 'bin' directory after the work done.
```sh
docker run -it \
    -v $(pwd)/:/proto/ \ # our proto files is here
    -v $(pwd)/client/:/client_src/ \ # client source code
    -v $(pwd)/:/hostoutput \     # output path for binary files
    tnt_grpc_client_builder
```

Now, we can to start the client for your platform.
```sh
./client/bin/client_mac
```
