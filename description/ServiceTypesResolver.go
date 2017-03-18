package description

import (
	"godegen/reflect"

	"github.com/bradfitz/slice"
)

var includedGenerics = map[string]bool{
	"System.Nullable`1":                 false, // don't resolve int? for example, but process the "int" part
	"System.Collections.Generic.List`1": false, // include lists
}

type ServiceTypesResolver struct {
	assembly       *reflect.Assembly
	typesByName    map[string]reflect.Type
	typesFindOrder []reflect.Type
}

func NewServiceTypesResolver(sourceAssembly *reflect.Assembly) *ServiceTypesResolver {
	typesByName := map[string]reflect.Type{}
	return &ServiceTypesResolver{
		sourceAssembly,
		typesByName,
		make([]reflect.Type, 0, 1024),
	}
}

func (res *ServiceTypesResolver) Resolve(serviceType reflect.Type) []reflect.Type {
	res.innerResolve(serviceType)
	return res.outputTypesAsSlice()
}

func includeGeneric(generic *reflect.GenericType) bool {
	baseName := generic.BaseType.FullName()
	include := includedGenerics[baseName] == true
	return include
}

func (res *ServiceTypesResolver) innerResolve(targetType reflect.Type) {
	include := true

	// exclude nil types
	if targetType == nil {
		return
	}

	// exclude built-in types
	if _, isBuiltIn := targetType.(*reflect.BuiltInType); isBuiltIn {
		return
	}

	// exclude generics, except specific types
	if generic, isGeneric := targetType.(*reflect.GenericType); isGeneric {
		if !includeGeneric(generic) {
			include = false
		}
	}

	// exclude arrays
	if _, isArray := targetType.(*reflect.ArrayType); isArray {
		include = false
	}

	typeName := targetType.FullName()
	if res.typesByName[typeName] != nil {
		return
	}

	if include {
		res.typesByName[typeName] = targetType
		res.typesFindOrder = append(res.typesFindOrder, targetType)
	}

	// generic arguments
	if genericType, ok := targetType.(*reflect.GenericType); ok {
		for _, argType := range genericType.ArgumentTypes() {
			res.innerResolve(argType)
		}
	}

	// array value type
	if array, isArray := targetType.(*reflect.ArrayType); isArray {
		res.innerResolve(array.ValueType())
	}

	// methods
	methods := targetType.GetMethods()
	for _, method := range methods {
		returnType := method.ReturnType()
		res.innerResolve(returnType)
		for _, param := range method.Parameters() {
			res.innerResolve(param.Type())
		}
	}

	// fields
	fields := targetType.GetFields()
	for _, field := range fields {
		fieldType := field.Type()
		res.innerResolve(fieldType)
	}

	// properties
	properties := targetType.GetProperties()
	for _, prop := range properties {
		propType := prop.Type()
		res.innerResolve(propType)
	}
}

func (res *ServiceTypesResolver) outputTypesAsSlice() []reflect.Type {
	//return res.typesFindOrder

	values := make([]reflect.Type, len(res.typesByName))
	i := 0
	for _, value := range res.typesByName {
		values[i] = value
		i++
	}

	slice.Sort(values, func(i, j int) bool {
		typeI := values[i]
		typeJ := values[j]
		if typeI.FullName() < typeJ.FullName() {
			return true
		}
		return false
	})

	return values
}
