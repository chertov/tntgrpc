package main

import (
    "github.com/emicklei/proto"
    "sort"
    "strings"
)

type Field struct {
    *proto.NormalField
    RealName string
    TypePath string
    FType FieldType
}
type FieldsMap map[string]Field;
func (h FieldsMap)getSortedKeys() []string {
    names := make([]string, 0, len(h))
    for name := range h {
        names = append(names, name)
    }
    sort.Strings(names)
    return names
}

type MapField struct {
    *proto.MapField
    RealName string
    TypePath string
    FType FieldType
}
type MapFieldsMap map[string]MapField;
func (h MapFieldsMap)getSortedKeys() []string {
    names := make([]string, 0, len(h))
    for name := range h {
        names = append(names, name)
    }
    sort.Strings(names)
    return names
}

type Message struct {
    *proto.Message
    Package string
    Path string
    Fields FieldsMap
    MapFields MapFieldsMap
}
type MessagesMap map[string]Message;
func (h MessagesMap)getSortedKeys() []string {
    names := make([]string, 0, len(h))
    for name := range h {
        names = append(names, name)
    }
    sort.Strings(names)
    return names
}

func NewMessage(ptr *proto.Message, pkgName string, path string, pd *ProtoData, pad string) Message {
    message := Message{}
    message.Message = ptr
    message.Package = pkgName
    message.Path = path + "." + message.Name
    message.Fields = make(FieldsMap)
    message.MapFields = make(MapFieldsMap)
    log_(pad, "Message", message.Path)
    pad += "    "
    for _, element := range message.Elements {
        switch v := element.(type) {
        case *proto.Enum:
            NewEnum(v, pkgName, message.Path, pd, pad)
        case *proto.Message:
            NewMessage(v, pkgName, message.Path, pd, pad)
        case *proto.NormalField:
            typePath := fixType(message.Path, v.Type)
            field := Field{v, v.Name, typePath, TSimple}
            field.Name = strings.ToLower(field.Name)
            log_(pad, "NormalField: ", field.TypePath, field.Name)
            message.Fields[field.Name] = field
        case *proto.MapField:
            if !isSimpleType(v.KeyType) {
                panic("Key type isn't simple type")
            }
            typePath := fixType(message.Path, v.Type)
            field := MapField{v, v.Name,typePath, TMap}
            field.Name = strings.ToLower(field.Name)
            log_(pad, "MapField: ", field.TypePath, field.Name)
            message.MapFields[field.Name] = field
        default:
        }
    }
    pd.Messages[message.Path] = message
    return message
}


func genCppToLuaForType(message Message, messages MessagesMap, header bool) string {
    if isSimpleType(message.Path) {
        return ""
    }
    cpp := "void tolua(lua_State *L, const std::string &key, const " + getCppType(message.Path) + " &obj)"
    if header {
        cpp += ";\n"
        return cpp
    }
    cpp += " {\n"
    cpp += "    if (key.size() > 0) lua_pushstring(L, key.c_str());\n"
    cpp += "    lua_newtable(L);\n"
    for _, fieldKey := range message.Fields.getSortedKeys() {
        field := message.Fields[fieldKey]
        if field.NormalField.Repeated {
            cpp += toLuaRepeatedField(field)
        } else {
            if field.FType == TSimple {
                cpp += "    tolua_" + field.TypePath + "(L, \"" + field.RealName + "\", obj." + field.Name + "()); // " + field.TypePath + "\n"
            } else if field.FType == TMap {
                // cpp += "    tolua(L, \"" + field.Name + "\", obj." + field.Name + "()); // " + field.TypePath + "\n"
            }  else if field.FType == TMessage {
                cpp += "    tolua(L, \"" + field.RealName + "\", obj." + field.Name + "()); // " + field.TypePath + "\n"
            }
        }
    }
    for _, fieldKey := range message.MapFields.getSortedKeys() {
        field := message.MapFields[fieldKey]
        cpp += toLuaMapField(field)
    }
    cpp += "    if (key.size() > 0) lua_settable( L, -3 );\n"
    cpp += "}\n"
    return cpp
}