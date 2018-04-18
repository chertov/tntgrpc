#pragma once

#ifdef __cplusplus
extern "C" {
#endif

#include <lua.h>
#include <lualib.h>
#include <lauxlib.h>

#ifndef GRPCLIBNAME
#define GRPCLIBNAME grpc
#endif

#define CONCAT1(a, b) a ## b
#define CONCAT(a, b) CONCAT1(a, b)
#define libopen_GRPC CONCAT(luaopen_, GRPCLIBNAME)

/* exported function */
LUA_API int libopen_GRPC(lua_State *L);

#ifdef __cplusplus
}
#endif
