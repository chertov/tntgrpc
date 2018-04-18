package main

import "strings"

func toLuaRepeatedField(field Field) string {
    if !field.Repeated {
        return ""
    }
    cpp := ""
    cpp += "// ============== '" + field.RealName + "' repeated[" + field.TypePath + "] ============== //\n"
    cpp += "lua_pushstring(L, \"" + field.RealName + "\");\n"
    cpp += "lua_newtable(L);\n"
    cpp += "for (unsigned int index = 0; index < obj." + field.Name + "_size(); ++index) {\n"
    cpp += "    lua_pushinteger(L, index); // push array index key\n"
    if isSimpleType(field.TypePath) {
        cpp += "    tolua_" + field.TypePath + "(L, obj." + field.Name + "(index)); // push value\n"
    } else {
        cpp += "    tolua(L, \"\", obj." + field.Name + "(index)); // push value\n"
    }
    cpp += "    lua_settable(L, -3);\n"
    cpp += "}\n"
    cpp += "lua_settable(L, -3);\n"
    cpp += "// ============== '" + field.RealName + "' repeated[" + field.TypePath + "] ============== //\n"
    return pad(0, cpp)
}

func fromLuaRepeatedField(field Field) string {
    cpp := `
            // Try to parse repeated of '` + field.TypePath + `' elements
            grpcTry(
                if (!lua_istable(L, -1)) return grpc_errors::RepeatedIsNotLuaTable::New("` + field.TypePath + `", getTypeName(lua_type(L, -1)))->run(SRC_POS);
                lua_pushnil(L); int64_t idx = 1;
                while(lua_next(L, -2) != 0) {
                    std::shared_ptr<void> _(nullptr, [&L](...){ lua_pop(L, 1); });
                    if (lua_type(L, -2) != LUA_TNUMBER) { return grpc_errors::RepeatedIndexIncorrectType::New("` + field.TypePath + `", idx, getTypeName(lua_type(L, -2)))->run(SRC_POS); }
                    double d = lua_tonumber(L, -2);
                    if (std::trunc(d) != d) { return grpc_errors::RepeatedIndexIsntInt::New("` + field.TypePath + `", idx, d)->run(SRC_POS); }
                    int64_t key = static_cast<int64_t>(d);
                    if(key != idx) { return grpc_errors::RepeatedIndexIncorrectOrder::New("` + field.TypePath + `", idx, key)->run(SRC_POS); }
        `
    cppType := getCppType(field.TypePath)
    if field.FType == TSimple {
        cpp += "            " + cppType + " val; grpcIfErr(fromlua_" + field.TypePath + "(L,val), grpc_errors::RepeatedValueIncorrectType::New(\"" + field.TypePath + "\", idx, getTypeName(lua_type(L, -1)))); obj->add_" + field.Name + "(val);\n";
    } else if field.FType == TEnum {
        cpp += "            " + cppType + " val; grpcIfErr(fromlua(L,fieldName.c_str(),val), grpc_errors::RepeatedValueIncorrectType::New(\"" + field.TypePath + "\", idx, getTypeName(lua_type(L, -1)))); obj->add_" + field.Name + "(val);\n";
    } else if field.FType == TMessage {
        cpp += "            " + "auto ptr = obj->add_" + strings.ToLower(field.Name) + "(); grpcIfErr(fromlua(L,\"\",ptr), grpc_errors::RepeatedValueIncorrectType::New(\"" + field.TypePath + "\", idx, getTypeName(lua_type(L, -1))));\n"
    }
    cpp += `                    idx++;
                }
            , grpc_errors::RepeatedParseFailure::New("` + field.TypePath + `"))
`
    return cpp
}