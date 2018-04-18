#include "luatools.hpp"
#include <string>
#include <iostream>
#include <cstdint>

#ifdef TARANTOOL_GRPC
    #include <tarantool/module.h>
#endif

#include "errors.hpp"

lua_State *L;

lua_State *getLuaState() {
#ifdef TARANTOOL_GRPC
    return luaT_state();
#else
    return L;
#endif
}
void setLuaState(lua_State *l) { L = l; }

GRPCError loadProc(lua_State *L, const std::string &service, const std::string &procedure) {
    lua_getglobal(L, service.c_str());
    if (!lua_istable(L, -1)) { return grpc_errors::ServiceNotFound::New(service)->run(SRC_POS); }
    lua_getfield(L, -1, procedure.c_str());
    if (!lua_isfunction(L, -1)) { return grpc_errors::ProcedureNotFound::New(service, procedure)->run(SRC_POS); }
    return NOERR;
}
GRPCError callProc(lua_State *L, const std::string &service, const std::string &procedure) {
    int top = lua_gettop(L);
    if (lua_pcall(L, 2, 2, 0) != 0) {
        std::stringstream ss; ss << "error running function: " << lua_tostring(L, -1);
        return grpc_errors::CallProcedureError::New(service, procedure, ss.str())->run(SRC_POS);;
    }
    return NOERR;
}


void tolua_double(lua_State *L, const char* key, const double &v) { lua_pushstring(L, key); lua_pushnumber(L, v); lua_settable(L, -3); }
void tolua_float(lua_State *L, const char* key, const float &v) { lua_pushstring(L, key); lua_pushnumber(L, v); lua_settable(L, -3); }
void tolua_int32(lua_State *L, const char* key, const int32_t &v) { lua_pushstring(L, key); lua_pushinteger(L, v); lua_settable(L, -3); }
void tolua_int64(lua_State *L, const char* key, const int64_t &v) { lua_pushstring(L, key); lua_pushinteger(L, v); lua_settable(L, -3); }
void tolua_uint32(lua_State *L, const char* key, const uint32_t &v) { lua_pushstring(L, key); lua_pushinteger(L, v); lua_settable(L, -3); }
void tolua_uint64(lua_State *L, const char* key, const uint64_t &v) { lua_pushstring(L, key); lua_pushinteger(L, v); lua_settable(L, -3); }
void tolua_sint32(lua_State *L, const char* key, const int32_t &v) { lua_pushstring(L, key); lua_pushinteger(L, v); lua_settable(L, -3); }
void tolua_sint64(lua_State *L, const char* key, const int64_t &v) { lua_pushstring(L, key); lua_pushinteger(L, v); lua_settable(L, -3); }
void tolua_fixed32(lua_State *L, const char* key, const uint32_t &v) { lua_pushstring(L, key); lua_pushinteger(L, v); lua_settable(L, -3); }
void tolua_fixed64(lua_State *L, const char* key, const uint64_t &v) { lua_pushstring(L, key); lua_pushinteger(L, v); lua_settable(L, -3); }
void tolua_sfixed32(lua_State *L, const char* key, const int32_t &v) { lua_pushstring(L, key); lua_pushinteger(L, v); lua_settable(L, -3); }
void tolua_sfixed64(lua_State *L, const char* key, const int64_t &v) { lua_pushstring(L, key); lua_pushinteger(L, v); lua_settable(L, -3); }
void tolua_bool(lua_State *L, const char* key, const bool &v) { lua_pushstring(L, key); lua_pushboolean(L, v); lua_settable(L, -3); }
void tolua_string(lua_State *L, const char* key, const std::string &v) { lua_pushstring(L, key); lua_pushlstring(L, v.data(), v.size()); lua_settable(L, -3); }
void tolua_bytes(lua_State *L, const char* key, const std::string &v) { lua_pushstring(L, key); lua_pushlstring(L, v.data(), v.size()); lua_settable(L, -3); }

void tolua_double(lua_State *L, const double &v) { lua_pushnumber(L, v); }
void tolua_float(lua_State *L, const float &v) { lua_pushnumber(L, v); }
void tolua_int32(lua_State *L, const int32_t &v) { lua_pushinteger(L, v); }
void tolua_int64(lua_State *L, const int64_t &v) { lua_pushinteger(L, v); }
void tolua_uint32(lua_State *L, const uint32_t &v) { lua_pushinteger(L, v); }
void tolua_uint64(lua_State *L, const uint64_t &v) { lua_pushinteger(L, v); }
void tolua_sint32(lua_State *L, const int32_t &v) { lua_pushinteger(L, v); }
void tolua_sint64(lua_State *L, const int64_t &v) { lua_pushinteger(L, v); }
void tolua_fixed32(lua_State *L, const uint32_t &v) { lua_pushinteger(L, v); }
void tolua_fixed64(lua_State *L, const uint64_t &v) { lua_pushinteger(L, v); }
void tolua_sfixed32(lua_State *L, const int32_t &v) { lua_pushinteger(L, v); }
void tolua_sfixed64(lua_State *L, const int64_t &v) { lua_pushinteger(L, v); }
void tolua_bool(lua_State *L, const bool &v) { lua_pushboolean(L, v); }
void tolua_string(lua_State *L, const std::string &v) { lua_pushlstring(L, v.data(), v.size()); }
void tolua_bytes(lua_State *L, const std::string &v) { lua_pushlstring(L, v.data(), v.size()); }

std::string getTypeName(int type) {
    switch (type) {
        case LUA_TNONE: return "LUA_TNONE";
        case LUA_TNIL: return "LUA_TNIL";
        case LUA_TBOOLEAN: return "LUA_TBOOLEAN";
        case LUA_TLIGHTUSERDATA: return "LUA_TLIGHTUSERDATA";
        case LUA_TNUMBER: return "LUA_TNUMBER";
        case LUA_TSTRING: return "LUA_TSTRING";
        case LUA_TTABLE: return "LUA_TTABLE";
        case LUA_TFUNCTION: return "LUA_TFUNCTION";
        case LUA_TUSERDATA: return "LUA_TUSERDATA";
        case LUA_TTHREAD: return "LUA_TTHREAD";
        // case LUA_NUMTAGS: return "LUA_NUMTAGS";
        default: return "Unknown";
    }
}

#define FROMLUA_CHECKERR(objtype) \
    int t = lua_type(L, -1+offset); \
    if (type != nullptr) *type = t; \
    if (t != objtype) { \
        return grpc_errors::IncorrectType::New(#objtype, getTypeName(lua_type(L, -1+offset)).c_str())->run(SRC_POS); \
    }
#define FROMLUA(lua_to) v = lua_to(L, -1+offset); return NOERR;

GRPCError fromlua_double(lua_State *L, double &v, const int offset, int *type) { FROMLUA_CHECKERR(LUA_TNUMBER); FROMLUA(lua_tonumber); }
GRPCError fromlua_float(lua_State *L, float &v, const int offset, int *type) { FROMLUA_CHECKERR(LUA_TNUMBER); FROMLUA(lua_tonumber); }
GRPCError fromlua_int32(lua_State *L, int32_t &v, const int offset, int *type) { FROMLUA_CHECKERR(LUA_TNUMBER); FROMLUA(lua_tointeger); }
GRPCError fromlua_int64(lua_State *L, int64_t &v, const int offset, int *type) { FROMLUA_CHECKERR(LUA_TNUMBER); FROMLUA(lua_tointeger); }
GRPCError fromlua_uint32(lua_State *L, uint32_t &v, const int offset, int *type) { FROMLUA_CHECKERR(LUA_TNUMBER); FROMLUA(lua_tointeger); }
GRPCError fromlua_uint64(lua_State *L, uint64_t &v, const int offset, int *type) { FROMLUA_CHECKERR(LUA_TNUMBER); FROMLUA(lua_tointeger); }
GRPCError fromlua_sint32(lua_State *L, int32_t &v, const int offset, int *type) { FROMLUA_CHECKERR(LUA_TNUMBER); FROMLUA(lua_tointeger); }
GRPCError fromlua_sint64(lua_State *L, int64_t &v, const int offset, int *type) { FROMLUA_CHECKERR(LUA_TNUMBER); FROMLUA(lua_tointeger); }
GRPCError fromlua_fixed32(lua_State *L, uint32_t &v, const int offset, int *type) { FROMLUA_CHECKERR(LUA_TNUMBER); FROMLUA(lua_tointeger); }
GRPCError fromlua_fixed64(lua_State *L, uint64_t &v, const int offset, int *type) { FROMLUA_CHECKERR(LUA_TNUMBER); FROMLUA(lua_tointeger); }
GRPCError fromlua_sfixed32(lua_State *L, int32_t &v, const int offset, int *type) { FROMLUA_CHECKERR(LUA_TNUMBER); FROMLUA(lua_tointeger); }
GRPCError fromlua_sfixed64(lua_State *L, int64_t &v, const int offset, int *type) { FROMLUA_CHECKERR(LUA_TNUMBER); FROMLUA(lua_tointeger); }
GRPCError fromlua_bool(lua_State *L, bool &v, const int offset, int *type) { FROMLUA_CHECKERR(LUA_TBOOLEAN); FROMLUA(lua_toboolean); }
GRPCError fromlua_string(lua_State *L, std::string &v, const int offset, int *type) { FROMLUA_CHECKERR(LUA_TSTRING); FROMLUA(lua_tostring); }
GRPCError fromlua_bytes(lua_State *L, std::string &v, const int offset, int *type) { FROMLUA_CHECKERR(LUA_TSTRING); FROMLUA(lua_tostring); }
