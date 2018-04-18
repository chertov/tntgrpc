package main

import "strings"

func toLuaMapField(field MapField) string {
    cpp := ""
    cpp += "// ============== '" + field.RealName + "' map<" + field.KeyType + ", " + field.TypePath + "> ============== //\n"
    cpp += "lua_pushstring(L, \"" + field.RealName + "\");\n"
    cpp += "lua_newtable(L);\n"
    cpp += "for (auto it = obj." + field.Name + "().begin(); it != obj." + field.Name + "().end(); it++ ) {\n"
    cpp += "    tolua_" + field.KeyType + "(L, it->first); // push key\n"
    if isSimpleType(field.TypePath) {
        cpp += "    tolua_" + field.TypePath + "(L, it->second); // push value\n"
    } else {
        cpp += "    tolua(L, \"\", it->second); // push value\n"
    }
    cpp += "    lua_settable(L, -3);\n"
    cpp += "}\n"
    cpp += "lua_settable(L, -3);\n"
    cpp += "// ============== '" + field.RealName + "' map<" + field.KeyType + ", " + field.TypePath + "> ============== //\n"
    return pad(0, cpp)
}

func fromLuaMapField(field MapField) string {
    keyCppType := getCppType(field.KeyType)
    cpp :=
`           // Try to parse map<` + field.KeyType + `,` + field.TypePath + `>
            if (!lua_istable(L, -1)) return grpc_errors::MapNotLuaTable::New("` + field.KeyType + `", "` + field.TypePath + `", getTypeName(lua_type(L, -1)))->run(SRC_POS);
            lua_pushnil(L);
            auto mapfield = obj->mutable_` + strings.ToLower(field.Name) + `();
            while(lua_next(L, -2) != 0) {
                std::shared_ptr<void> _(nullptr, [&L](...){ lua_pop(L, 1); });
                ` + keyCppType + ` key; int luatype; grpcIfErr(fromlua_` + field.KeyType + `(L,key,-1,&luatype), grpc_errors::MapElementKeyIncorrectType::New("` + field.KeyType + `","` + field.TypePath + `", lua_tostring(L, -2), getTypeName(luatype)));
`
    cppType := getCppType(field.TypePath)
    if field.FType == TSimple {
        cpp += "                " + cppType + " val; grpcIfErr(fromlua_" + field.TypePath + "(L,val), grpc_errors::MapElementValueIncorrectType::New(\"" + field.KeyType + "\",\"" + field.TypePath + "\", key, getTypeName(lua_type(L, -1)) ));\n";
    } else if field.FType == TEnum {
        cpp += "                " + cppType + " val; grpcIfErr(fromlua(L,\"\",val), grpc_errors::MapElementValueIncorrectType::New(\"" + field.KeyType + "\",\"" + field.TypePath + "\", key, getTypeName(lua_type(L, -1)) ));\n"
    } else {
        cpp += "                " + cppType + " val; grpcIfErr(fromlua(L,\"\",&val), grpc_errors::MapElementValueIncorrectType::New(\"" + field.KeyType + "\",\"" + field.TypePath + "\", key, getTypeName(lua_type(L, -1)) ));\n"
    }
    cpp += "                (*mapfield)[key] = val;\n"
    cpp += "            }"
    return cpp
}
