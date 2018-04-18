package main

import (
    "strings"
)

func isSimpleType(name string) bool {
    if  name == "double" ||
        name == "float" ||
        name == "int32" ||
        name == "int64" ||
        name == "uint32" ||
        name == "uint64" ||
        name == "sint32" ||
        name == "sint64" ||
        name == "fixed32" ||
        name == "fixed64" ||
        name == "sfixed32" ||
        name == "sfixed64" ||
        name == "bool" ||
        name == "string" ||
        name == "bytes" {
        return true
    }
    return false
}

func getCppType(name string) string {
    switch name {
    case "double": return "double"
    case "float": return "float"
    case "int32": return "int32_t"
    case "int64": return "int64_t"
    case "uint32": return "uint32_t"
    case "uint64": return "uint64_t"
    case "sint32": return "int32_t"
    case "sint64": return "int64_t"
    case "fixed32": return "uint32_t"
    case "fixed64": return "uint64_t"
    case "sfixed32": return "int32_t"
    case "sfixed64": return "int64_t"
    case "bool": return "bool"
    case "string": return "std::string"
    case "bytes": return "std::string"
    }
    return strings.Join(strings.Split(name, "."), "::")
}
func fixType(path, name string) string {
    if isSimpleType(name) {
        return name
    }
    return path + "." + name
}

func padStr(padCnt int) string {
    padstr := ""
    for i := 0; i < padCnt; i++ {
        padstr = padstr + "    "
    }
    return padstr
}

func pad(cnt int, str string) string {
    strs := make([]string, 0)
    for _, str := range strings.Split(str, "\n") {
        if len(str) > 0 {
            strs = append(strs, padStr(cnt) + str)
        } else {
            strs = append(strs, str)
        }
    }
    return strings.Join(strs, "\n")
}
func log_(v ...interface{}) {
    // log.Println(v...)
}