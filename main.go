package main

import (
	"fmt"
	"godegen/codegen"
	"godegen/reflect"
	"os"
	"strings"
)

const (
	SERVICE_PREFIX     = "nz.co.LanguagePerfect.Service"
	SERVICE_PREFIX_LEN = len(SERVICE_PREFIX)
)

var serviceNames = []string{
	"nz.co.LanguagePerfect.Services.Portals.ControlPanel.ContentPortal",
	"nz.co.LanguagePerfect.Services.Portals.ControlPanel.CRMPortal",
	"nz.co.LanguagePerfect.Services.Portals.ControlPanel.LanguageDataPortal",
	"nz.co.LanguagePerfect.Services.Portals.ControlPanel.QualityPortal",
	"nz.co.LanguagePerfect.Services.Portals.ControlPanel.UserTasksPortal",
}

func main() {
	describer := codegen.NewServiceDescriber(`C:\WF\LP\server\EBS_Deployment\bin`, `Classes`, typeMapper, namespaceMapper)
	for _, serviceName := range serviceNames {
		descr, _ := describer.Describe(serviceName)
		var gen *codegen.Generator
		var err error
		config := &codegen.GeneratorConfig{
			TemplatesPath: `.\templates\typescript`,
			FileExtension: ".ts",
		}
		if gen, err = codegen.NewGenerator(`c:\dooschmonkey`, config); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err = gen.OutputServiceDescription(descr); err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
	}

	// fmt.Print(descr.JSON())
}

var typescriptTypeMap = map[string]string{
	"System.Boolean":                   "boolean",
	"System.Byte":                      "number",
	"System.UInt16":                    "number",
	"System.UInt32":                    "number",
	"System.UInt64":                    "number",
	"System.SByte":                     "number",
	"System.Int16":                     "number",
	"System.Int32":                     "number",
	"System.Int64":                     "number",
	"System.Decimal":                   "number",
	"System.Double":                    "number",
	"System.DateTime":                  "Date",
	"System.Nullable<System.Byte>":     "number | null",
	"System.Nullable<System.SByte>":    "number | null",
	"System.Nullable<System.Int16>":    "number | null",
	"System.Nullable<System.Int32>":    "number | null",
	"System.Nullable<System.Int64>":    "number | null",
	"System.Nullable<System.Decimal>":  "number | null",
	"System.Nullable<System.Double>":   "number | null",
	"System.Nullable<System.DateTime>": "Date | null",
	"System.String":                    "string",
	"System.Object":                    "any",
	"System.ValueType":                 "any",
	"System.Void":                      "void",
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
		if isBuiltIn(elemType) {
			return typeMapper(elemType) + "[]"
		}
		return "Array<" + typeMapper(elemType) + ">"
	}

	return cleanedFullName
}

func isBuiltIn(typ reflect.Type) bool {
	_, isBuiltIn := typ.(*reflect.BuiltInType)
	return isBuiltIn
}

func isCollectionType(typ reflect.Type) (reflect.Type, bool) {
	if array, isArray := typ.(*reflect.ArrayType); isArray {
		return array.ValueType(), true
	}

	if generic, isGeneric := typ.(*reflect.GenericType); isGeneric {
		if generic.TypeBase.FullName() == "System.Collections.Generic.List`1" {
			return generic.ArgumentTypes()[0], true
		}
	}

	return nil, false
}
