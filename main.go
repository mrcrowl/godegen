package main

import (
	"fmt"
	"godegen/description"
	"godegen/reflect"
	"strings"
)

const (
	SERVICE_PREFIX     = "nz.co.LanguagePerfect.Service"
	SERVICE_PREFIX_LEN = len(SERVICE_PREFIX)
)

func main() {
	describer := description.NewServiceDescriber(`C:\WF\LP\server\EBS_Deployment\bin`, `Classes`, typeMapper, namespaceMapper)
	descr, _ := describer.Describe("nz.co.LanguagePerfect.Services.Portals.ControlPanel.LanguageDataPortal")
	fmt.Print(descr.JSON())
}

var typescriptTypeMap = map[string]string{
	"System.Byte":                      "number",
	"System.UInt16":                    "number",
	"System.UInt32":                    "number",
	"System.UInt64":                    "number",
	"System.SByte":                     "number",
	"System.Int16":                     "number",
	"System.Int32":                     "number",
	"System.Int64":                     "number",
	"System.String":                    "string",
	"System.Boolean":                   "bool",
	"System.DateTime":                  "Date",
	"System.Nullable<System.Byte>":     "(number | null)",
	"System.Nullable<System.SByte>":    "(number | null)",
	"System.Nullable<System.Int16>":    "(number | null)",
	"System.Nullable<System.Int32>":    "(number | null)",
	"System.Nullable<System.Int64>":    "(number | null)",
	"System.Nullable<System.DateTime>": "(number | null)",
}

func namespaceMapper(namespace string) string {
	if strings.HasPrefix(namespace, SERVICE_PREFIX) {
		return "service" + namespace[SERVICE_PREFIX_LEN:]
	}

	return namespace
}

func typeMapper(typ reflect.Type) string {
	fullname := typ.FullName()
	if mappedName, found := typescriptTypeMap[fullname]; found {
		return mappedName
	}

	cleanedFullName := namespaceMapper(fullname)
	if elemType, isCollection := isCollectionType(typ); isCollection {
		return typeMapper(elemType) + "[]"
	}

	return cleanedFullName
}

func isCollectionType(typ reflect.Type) (reflect.Type, bool) {
	if array, isArray := typ.(*reflect.ArrayType); isArray {
		return array.ValueType(), true
	}

	if generic, isGeneric := typ.(*reflect.GenericType); isGeneric {
		if generic.BaseType.FullName() == "System.Collections.Generic.List`1" {
			return generic.ArgumentTypes()[0], true
		}
	}

	return nil, false
}
