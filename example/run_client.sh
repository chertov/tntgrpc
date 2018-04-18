
# Now, we should build GRPC client for Tarantool GRPC service.
# Let's build it with docker. It's supereasy!

# We have Dockerfile with build tools, grpc and golang compilers. Let's build it.
docker build -t tnt_grpc_client_builder -f ./client/Dockerfile ./

# Let's start 'tnt_grpc_client_builder' container.
# We need to share *.proto files and Goland source code for it.
# Now we have the clients for Mac, Linux and Windows in 'bin' directory.
docker run -it -v $(pwd)/:/hostoutput \
    -v $(pwd)/client/:/client_src/ \
    -v $(pwd)/:/proto/ \
    tnt_grpc_client_builder

case "$OSTYPE" in
    darwin*)  ./client/bin/client_mac ;;
    linux*)   ./client/bin/client_linux ;;
    msys*)    ./client/bin/client.exe ;;
    *)        ./client/bin/client_linux ;;
esac
