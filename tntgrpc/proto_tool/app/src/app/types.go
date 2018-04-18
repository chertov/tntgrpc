package main

type FieldType int

const (
    TSimple FieldType = iota
    TMessage
    TEnum
    TMap
    TOneOf
    TAny
)

type CppClass struct {
    NameSpace string
    Name string
}
