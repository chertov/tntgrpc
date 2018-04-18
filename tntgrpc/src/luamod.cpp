#pragma once
#include <tarantool/module.h>
#include "luamod.h"
#include "grpc.hpp"

static int grpc_start(struct lua_State *L) {
    int argc = lua_gettop(L);
    std::string server_address = "0.0.0.0:50051";
    bool sync_mode = false;
    if (argc > 1 && lua_type(L, 1-argc) == LUA_TSTRING) server_address = lua_tostring(L, 1-argc);
    if (argc > 2 && lua_type(L, 2-argc) == LUA_TBOOLEAN) sync_mode = lua_toboolean(L, 2-argc);
    grpc_start(server_address, sync_mode);
    return 0;
}

/* exported function */
int libopen_GRPC(lua_State *L) {
    lua_newtable(L);
    static const struct luaL_Reg meta [] = {
        {"start", grpc_start},
        {NULL, NULL}
    };
    luaL_register(L, NULL, meta);
    return 1;
}
