package main

import (
    "github.com/emicklei/proto"
    "sort"
)

type EnumValue struct {
    Name string
    Integer int
}
type EnumValuesMap map[string]EnumValue;
func (h EnumValuesMap)getSortedKeys() []string {
    names := make([]string, 0, len(h))
    for name := range h {
        names = append(names, name)
    }
    sort.Strings(names)
    return names
}

type Enum struct {
    *proto.Enum
    Package string
    Path string
    Values EnumValuesMap
}
type EnumsMap map[string]Enum;
func (h EnumsMap)getSortedKeys() []string {
    names := make([]string, 0, len(h))
    for name := range h {
        names = append(names, name)
    }
    sort.Strings(names)
    return names
}

func NewEnum(v *proto.Enum, packageName string, path string, pd *ProtoData, pad string) {
    enum := Enum{}
    enum.Enum = v
    enum.Package = packageName
    enum.Path = path + "." + enum.Name
    enum.Values = make(map[string]EnumValue)
    log_(pad + "    ", "Enum ", enum.Path)
    for _, element := range v.Elements {
        value, ok := element.(*proto.EnumField)
        if ok {
            log_(pad + "    " + "    ", value.Name, value.Integer)
            enum.Values[value.Name] = EnumValue{value.Name, value.Integer}
        }
    }
    pd.Enums[enum.Path] = enum
}

func (f EnumsMap)toLua(header bool) string {
    cpp := ""
    for _, path := range f.getSortedKeys() {
        enum := f[path]
        cppType := getCppType(enum.Path)
        cpp = "GRPCError tolua(lua_State *L, const std::string &key, const " + cppType + " &obj)"
        if header {
            cpp += ";\n"
            return cpp
        }
        cpp += " {\n"
        cpp += "    std::string value;\n"
        cpp += "    switch (obj) {\n"
        for _, name := range enum.Values.getSortedKeys() {
            // enumValue := enum.Values[name]
            cpp += "    case " + cppType + "::" + name + ": { value = \"" + name + "\"; break; }\n"
        }
        cpp += "    default: { return grpc_errors::toLuaEnumValueIsIncorrect::New(\"" + enum.Path + "\", obj)->run(SRC_POS); }\n"
        cpp += "    }\n"
        cpp += "    tolua_string(L, key.c_str(), value);\n"
        cpp += "    return NOERR;\n"
        cpp += "}\n"
    }
    return cpp
}

func (f EnumsMap)fromLua(header bool) string {
    cpp := ""
    for _, path := range f.getSortedKeys() {
        enum := f[path]
        cppType := getCppType(enum.Path)
        cpp = "GRPCError fromlua(lua_State *L, const std::string &key, " + cppType + " &obj)"
        if header {
            cpp += ";\n"
            return cpp
        }
        cpp += " {\n"
        cpp += "    GRPCError err = NOERR;\n"
        cpp += "    std::string str; err = fromlua_string(L, str);\n"
        cpp += "    if (err->noerr()) {\n        "
        for _, name := range enum.Values.getSortedKeys() {
            cpp += "if (str.compare(\"" + name + "\") == 0) { obj = " + cppType + "::" + name + "; return NOERR; }\n        else "
        }
        cpp += "return grpc_errors::EnumValueIsntIncorrect::New(\"" + enum.Path + "\", str)->run(SRC_POS);\n"
        cpp += "    } else {\n"
        cpp += "        int32_t value; int _type; grpcIfErr(fromlua_int32(L, value, 0, &_type), grpc_errors::EnumIsNotStringOrInt::New(\"" + enum.Path + "\", getTypeName(_type)));\n        "
        for _, name := range enum.Values.getSortedKeys() {
            cpp += "if (value == static_cast<int32_t>(" + cppType + "::" + name + ")) { obj = " + cppType + "::" + name + "; return NOERR; }\n        else "
        }
        cpp += "return grpc_errors::EnumValueIsntIncorrect::New(\"" + enum.Path + "\", value)->run(SRC_POS);\n"
        cpp += "    }\n"
        cpp += "}\n"
    }
    return cpp
}
