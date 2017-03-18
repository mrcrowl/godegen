package description

import (
	"fmt"
	"godegen/reflect"
)

type ServiceDescriber struct {
	assemblyFile *reflect.Assembly
}

func NewServiceDescriber(assemblyFilepath string, assemblyName string) *ServiceDescriber {
	loader := reflect.NewAssemblyLoader(assemblyFilepath)
	assemblyFile, _ := loader.Load(assemblyName)
	return &ServiceDescriber{assemblyFile}
}

func (res *ServiceDescriber) Describe(serviceTypeName string) (*ServiceDescription, error) {
	serviceType := res.assemblyFile.GetType(serviceTypeName)
	resolver := NewServiceTypesResolver(res.assemblyFile)
	resolvedTypes := resolver.Resolve(serviceType)
	for _, t := range resolvedTypes {
		fmt.Println(t.FullName() + ":")
		for _, field := range t.GetFields() {
			fieldType := field.Type()
			if collectionType, isCollection := isCollectionType(fieldType); isCollection {
				fmt.Println("\t" + field.Name() + ": " + collectionType.FullName() + "[]")
			} else {
				fmt.Println("\t" + field.Name() + ": " + field.Type().FullName())
			}
		}
	}
	return nil, nil
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
