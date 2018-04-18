package main

import (
    "os"

    "github.com/emicklei/proto"
    "io/ioutil"
    "strings"
    "errors"
    "flag"
    "fmt"
    "path/filepath"
)

func genCppFromLuaForType(message Message, messages MessagesMap, header bool) string {
    if isSimpleType(message.Path) {
        return ""
    }
    cpp := "GRPCError fromlua(lua_State *L, const std::string &key, " + getCppType(message.Path) + " *obj)"
    if header {
        cpp += ";\n"
        return cpp
    }
    cpp += " {\n"
    cpp += "    static const std::string typeName = \"" + message.Path + "\";\n"
    // cpp += "    grpcTry (\n"
    cpp += "    if (!lua_istable(L, -1)) return grpc_errors::MessageNotLuaTable::New(\"" + message.Path + "\", getTypeName(lua_type(L, -1)))->run(SRC_POS);\n"
    cpp += "    lua_pushnil(L);\n"
    cpp += "    while(lua_next(L, -2) != 0) {\n"
    cpp += "        std::shared_ptr<void> _(nullptr, [&key, &L](...){ lua_pop(L, 1); }); // defer lua_pop(L, 1);\n"
    cpp += "        std::string fieldName = lua_tostring(L, -2);\n"
    cpp += "        "
    for _, fieldKey := range message.Fields.getSortedKeys() {
        field := message.Fields[fieldKey]
        cppType := getCppType(field.TypePath)
        cpp += "if (fieldNameCompare(fieldName,\"" + strings.ToLower(field.Name) + "\")) {"
        if field.Repeated {
            cpp += "grpcTry(" + fromLuaRepeatedField(field) + ", grpc_errors::MessageFieldParseFailure::New(\"" + message.Path + "\", \"" + field.Name + "\", \"" + field.TypePath + "\"));"
        } else {
            if field.FType == TSimple {
                cpp += " " + cppType + " val; grpcIfErr(fromlua_" + field.TypePath + "(L, val), grpc_errors::MessageFieldParseFailure::New(\"" + message.Path + "\", \"" + field.Name + "\", \"" + field.TypePath + "\")); obj->set_" + field.Name + "(val);"
            } else if field.FType == TEnum {
                cpp += " " + cppType + " val; grpcIfErr(fromlua(L, \"" + field.Name + "\", val), grpc_errors::MessageFieldParseFailure::New(\"" + message.Path + "\", \"" + field.Name + "\", \"" + field.TypePath + "\")); obj->set_" + field.Name + "(val);"
            } else if field.FType == TMessage {
                cpp += " grpcIfErr(fromlua(L, \"" + field.Name + "\", obj->mutable_" + field.Name + "()), grpc_errors::MessageFieldParseFailure::New(\"" + message.Path + "\", \"" + field.Name + "\", \"" + field.TypePath + "\")); /* " + field.TypePath + " */"
            } else {
                panic("Unknown type of field")
            }
        }
        cpp += "        }\n        else "
    }
    for _, fieldKey := range message.MapFields.getSortedKeys() {
        field := message.MapFields[fieldKey]
        cpp += "if (fieldNameCompare(fieldName, \"" + strings.ToLower(field.Name) + "\")) {\n"
        cpp += "            grpcTry(\n" + fromLuaMapField(field) + "\n"
        cpp += "            , grpc_errors::MessageFieldParseFailure::New(\"" + message.Path + "\", \"" + field.Name + "\", \"" + field.TypePath + "\"));\n"
        cpp += "        }\n        else "
    }

    cpp += "{ return grpc_errors::MessageNotExistingField::New(\"" + message.Path + "\", fieldName)->run(SRC_POS); }\n"

    cpp += "    }\n"
    // cpp += "    , grpc_errors::MessageFieldParseFailure::New(\"" + message.Path + "\", key, ));\n"
    cpp += "    return NOERR;\n"
    cpp += "}\n"
    return cpp
}

func genCpp(pd ProtoData, output string) string {
    log_()
    log_()
    cpp := ""
    cpp += pd.Enums.toLua(true)
    cpp += pd.Enums.fromLua(true)
    for _, path := range pd.Messages.getSortedKeys() {
        message := pd.Messages[path]
        cpp += genCppToLuaForType(message, pd.Messages, true)
        cpp += genCppFromLuaForType(message, pd.Messages, true)
    }
    cpp += "\n"
    cpp += pd.Enums.toLua(false)
    cpp += pd.Enums.fromLua(false)
    for _, path := range pd.Messages.getSortedKeys() {
        message := pd.Messages[path]
        cpp += genCppToLuaForType(message, pd.Messages, false)
        cpp += genCppFromLuaForType(message, pd.Messages, false)
    }
    cpp += "\n"
    cpp += pd.Services.genCpp()
    return cpp
}

func findType(typepath string, fieldtype string, pb ProtoData) (string, FieldType, error) {
    originType := typepath
    if isSimpleType(typepath) {
        return typepath, TSimple, nil
    }
    _, ok := pb.Messages[typepath]
    if ok {
        return typepath, TMessage, nil
    }
    _, ok = pb.Enums[typepath]
    if ok {
        return typepath, TEnum, nil
    }
    typepath = strings.TrimSuffix(typepath, "." + fieldtype)
    log_(originType, fieldtype, typepath)
    for true {
        paths := strings.Split(typepath, ".")
        if len(paths) > 0 {
            i := len(paths) - 1
            paths = append(paths[:i], paths[i+1:]...)
        }
        typepath = strings.Join(paths, ".")
        p := fieldtype
        if len(typepath) > 0 {
            p = typepath + "." + p
        }
        log_("in for ", typepath, fieldtype, p)
        _, ok = pb.Messages[p]
        if ok {
            return p, TMessage, nil
        }
        _, ok = pb.Enums[p]
        if ok {
            return p, TEnum, nil
        }
        if len(paths) == 0 {
            break
        }
    }
    return "", TSimple, errors.New("Can't find type with name " + originType)
}

type ProtoData struct {
    Services ServicesMap
    Messages MessagesMap
    Enums EnumsMap
}
func NewProtoData() ProtoData {
    pd := ProtoData{}
    pd.Services = make(ServicesMap)
    pd.Messages = make(MessagesMap)
    pd.Enums = make(EnumsMap)
    return pd
}

func (pd *ProtoData) ParseProto(filepath string) {
    reader, err := os.Open(filepath)
    defer reader.Close()
    if err != nil {
        panic(err)
    }
    parser := proto.NewParser(reader)
    definition, _ := parser.Parse()
    pkgName := ""
    pad := ""
    for _, element := range definition.Elements {
        switch v := element.(type) {
        case *proto.Import:
            log_("proto.Import: ", v.Filename)
        case *proto.Option:
            log_("proto.Option: ", v.Name)
        case *proto.Package:
            log_("proto.packageName: ", v.Name)
            pkgName = v.Name
        case *proto.Service:
            NewService(v, pkgName, pd, pad+"    ")
        case *proto.Message:
            NewMessage(v, pkgName, pkgName, pd, pad+"    ")
        case *proto.Enum:
            NewEnum(v, pkgName, pkgName, pd, pad+"    ")
        default:
        }
    }
    log_()
    log_("package: ", pkgName)
    log_()
    log_()
}

func (pd *ProtoData) FixTypes(output string) {
    // fix type paths
    for _, path := range pd.Messages.getSortedKeys() {
        message := pd.Messages[path]
        log_(path, message.Name)
        for _, name := range message.Fields.getSortedKeys() {
            field := message.Fields[name]
            fixPath, fType, err := findType(field.TypePath, field.Type, *pd)
            if err != nil {
                panic(err)
            }
            field.FType = fType
            field.TypePath = fixPath
            log_("field: ", name, field.TypePath)
            message.Fields[name] = field
        }
        for _, name := range message.MapFields.getSortedKeys() {
            field := message.MapFields[name]
            fixPath, fType, err := findType(field.TypePath, field.Type, *pd)
            if err != nil {
                panic(err)
            }
            field.FType = fType
            field.TypePath = fixPath
            log_("field: ", name, field.TypePath)
            message.MapFields[name] = field
        }
    }

    cpp := ""
    protofiles := flag.Args()
    for _, protopath := range protofiles {
        cpp += "#include \"" + strings.TrimSuffix(filepath.Base(protopath), filepath.Ext(protopath)) + ".pb.h\"\n"
        cpp += "#include \"" + strings.TrimSuffix(filepath.Base(protopath), filepath.Ext(protopath)) + ".grpc.pb.h\"\n"
        cpp += "#include \"" + strings.TrimSuffix(filepath.Base(protopath), filepath.Ext(protopath)) + ".pb.cc\"\n"
        cpp += "#include \"" + strings.TrimSuffix(filepath.Base(protopath), filepath.Ext(protopath)) + ".grpc.pb.cc\"\n"
    }
    cpp += "\n"
    cpp += genCpp(*pd, output)
    err := ioutil.WriteFile(output, []byte(cpp), 0644)
    if err != nil {
        panic(err)
    }
}

func main() {
    pd := NewProtoData()
    outputPath := flag.String("outcpp", "./gen/grpc.gen.cpp", "C++ output path")
    // protoIpath := flag.String("I", "", "proto include paths")
    flag.Usage = func() {
        tool := filepath.Base(os.Args[0])
        fmt.Fprintf(os.Stderr, "Usage of %s:\n", tool)
        flag.PrintDefaults()
        fmt.Fprintf(os.Stderr, "examples:\n")
        fmt.Fprintf(os.Stderr, "    %s --outcpp=./gen/grpc.gen.cpp  helloworld.proto test.proto\n", tool)
        fmt.Fprintf(os.Stderr, "    %s helloworld.proto test.proto\n", tool)
    }
    flag.Parse()

    protofiles := flag.Args()
    if len(protofiles) == 0 {
        flag.Usage()
        os.Exit(1)
    }

    for _, filepath := range protofiles {
        pd.ParseProto(filepath)
    }
    pd.FixTypes(*outputPath)
}
