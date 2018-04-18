#pragma once

#include <string>

#ifdef __cplusplus
extern "C" {
#endif
#define LUA_COMPAT_MODULE
#include <lua.h>
#include <lualib.h>
#include <lauxlib.h>

#ifdef __cplusplus
}
#endif

#include "errors.hpp"

lua_State *getLuaState();
void setLuaState(lua_State *);
GRPCError loadProc(lua_State *L, const std::string &service, const std::string &procedure);
GRPCError callProc(lua_State *L, const std::string &service, const std::string &procedure);

std::string getTypeName(int type);

void tolua_double(lua_State *L, const char* key, const double &v); // double
void tolua_float(lua_State *L, const char* key, const float &v); // float
void tolua_int32(lua_State *L, const char* key, const int32_t &v); // int32
void tolua_int64(lua_State *L, const char* key, const int64_t &v); // int64
void tolua_uint32(lua_State *L, const char* key, const uint32_t &v); // uint32
void tolua_uint64(lua_State *L, const char* key, const uint64_t &v); // uint64
void tolua_sint32(lua_State *L, const char* key, const int32_t &v); // sint32
void tolua_sint64(lua_State *L, const char* key, const int64_t &v); // sint64
void tolua_fixed32(lua_State *L, const char* key, const uint32_t &v); // fixed32
void tolua_fixed64(lua_State *L, const char* key, const uint64_t &v); // fixed64
void tolua_sfixed32(lua_State *L, const char* key, const int32_t &v); // sfixed32
void tolua_sfixed64(lua_State *L, const char* key, const int64_t &v); // sfixed64
void tolua_bool(lua_State *L, const char* key, const bool &v); // bool
void tolua_string(lua_State *L, const char* key, const std::string &v); // string
void tolua_bytes(lua_State *L, const char* key, const std::string &v); // bytes


void tolua_double(lua_State *L, const double &v); // double
void tolua_float(lua_State *L, const float &v); // float
void tolua_int32(lua_State *L, const int32_t &v); // int32
void tolua_int64(lua_State *L, const int64_t &v); // int64
void tolua_uint32(lua_State *L, const uint32_t &v); // uint32
void tolua_uint64(lua_State *L, const uint64_t &v); // uint64
void tolua_sint32(lua_State *L, const int32_t &v); // sint32
void tolua_sint64(lua_State *L, const int64_t &v); // sint64
void tolua_fixed32(lua_State *L, const uint32_t &v); // fixed32
void tolua_fixed64(lua_State *L, const uint64_t &v); // fixed64
void tolua_sfixed32(lua_State *L, const int32_t &v); // sfixed32
void tolua_sfixed64(lua_State *L, const int64_t &v); // sfixed64
void tolua_bool(lua_State *L, const bool &v); // bool
void tolua_string(lua_State *L, const std::string &v); // string
void tolua_bytes(lua_State *L, const std::string &v); // bytes


GRPCError fromlua_double(lua_State *L, double &v, const int offset = 0, int *type = nullptr); // double
GRPCError fromlua_float(lua_State *L, float &v, const int offset = 0, int *type = nullptr); // float
GRPCError fromlua_int32(lua_State *L, int32_t &v, const int offset = 0, int *type = nullptr); // int32
GRPCError fromlua_int64(lua_State *L, int64_t &v, const int offset = 0, int *type = nullptr); // int64
GRPCError fromlua_uint32(lua_State *L, uint32_t &v, const int offset = 0, int *type = nullptr); // uint32
GRPCError fromlua_uint64(lua_State *L, uint64_t &v, const int offset = 0, int *type = nullptr); // uint64
GRPCError fromlua_sint32(lua_State *L, int32_t &v, const int offset = 0, int *type = nullptr); // sint32
GRPCError fromlua_sint64(lua_State *L, int64_t &v, const int offset = 0, int *type = nullptr); // sint64
GRPCError fromlua_fixed32(lua_State *L, uint32_t &v, const int offset = 0, int *type = nullptr); // fixed32
GRPCError fromlua_fixed64(lua_State *L, uint64_t &v, const int offset = 0, int *type = nullptr); // fixed64
GRPCError fromlua_sfixed32(lua_State *L, int32_t &v, const int offset = 0, int *type = nullptr); // sfixed32
GRPCError fromlua_sfixed64(lua_State *L, int64_t &v, const int offset = 0, int *type = nullptr); // sfixed64
GRPCError fromlua_bool(lua_State *L, bool &v, const int offset = 0, int *type = nullptr); // bool
GRPCError fromlua_string(lua_State *L, std::string &v, const int offset = 0, int *type = nullptr); // string
GRPCError fromlua_bytes(lua_State *L, std::string &v, const int offset = 0, int *type = nullptr); // bytes
