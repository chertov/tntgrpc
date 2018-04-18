cmake_minimum_required(VERSION 3.6)
project(grpctntso)

set(CMAKE_CXX_STANDARD 14)

include_directories(
    ./
    /usr/local/include/
)
link_directories(
    /usr/local/lib/
)

set(GRPCTNT_SOURCES
    src/errors.hpp
    src/errors.cpp
    src/luatools.hpp
    src/luatools.cpp
    src/luamod.h
    src/luamod.cpp
    src/grpc.hpp
    src/grpc.cpp
)
add_definitions(
    -DTARANTOOL_GRPC
    -DGRPCLIBNAME=$ENV{TNT_GRPC_NAME}
)
add_library(grpctntso SHARED ${GRPCTNT_SOURCES})
target_link_libraries(grpctntso
    libprotobuf.a
    libgrpc++.a
    libgrpc.a
)
set_target_properties(grpctntso
    PROPERTIES
    ARCHIVE_OUTPUT_DIRECTORY "/output/"
    LIBRARY_OUTPUT_DIRECTORY "/output/"
    RUNTIME_OUTPUT_DIRECTORY "/output/"
    PREFIX ""
    OUTPUT_NAME $ENV{TNT_GRPC_NAME}
)
