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
			fmt.Println("\t" + field.Name() + ": " + field.Type().FullName())
		}
	}
	return nil, nil
}
