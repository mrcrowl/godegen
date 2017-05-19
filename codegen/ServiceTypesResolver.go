package codegen

import (
	"godegen/reflect"

	"github.com/bradfitz/slice"
)

var genericsWhitelist = map[string]bool{ // false indicates don't include the generic itself
	"System.Nullable`1":                 false, // don't resolve int? for example, but process the "int" part
	"System.Collections.Generic.List`1": false, // include list elements
	"System.Threading.Tasks.Task`1":     false, // include task wrapped types
}

var excludedBaseTypes = map[string]bool{
	"System.Object":                                               true,
	"System.ValueType":                                            true,
	"nz.co.LanguagePerfect.Services.PortalsAsync.AsyncPortalBase": true,
	"System.Web.Services.WebService":                              true,
	"System.ComponentModel.MarshalByValueComponent":               true,
}

type ServiceTypesResolver struct {
	assembly       *reflect.Assembly
	serviceType    reflect.Type
	typesByName    map[string]reflect.Type
	typesFindOrder []reflect.Type
}

func ResolveServiceDependencyTypes(sourceAssembly *reflect.Assembly, serviceType reflect.Type) []reflect.Type {
	resolver := newServiceTypesResolver(sourceAssembly, serviceType)
	return resolver.resolve()
}

func newServiceTypesResolver(sourceAssembly *reflect.Assembly, serviceType reflect.Type) *ServiceTypesResolver {
	typesByName := map[string]reflect.Type{}
	return &ServiceTypesResolver{
		sourceAssembly,
		serviceType,
		typesByName,
		make([]reflect.Type, 0, 1024),
	}
}

func (res *ServiceTypesResolver) resolve() []reflect.Type {
	res.innerResolve(res.serviceType)
	return res.outputTypesAsSlice()
}

func includeGeneric(generic *reflect.GenericType) bool {
	baseName := generic.TypeBase.FullName()
	include := genericsWhitelist[baseName] == true
	return include
}

func (res *ServiceTypesResolver) innerResolve(targetType reflect.Type) {
	isServiceType := targetType == res.serviceType
	include := !isServiceType

	// exclude nil types
	if targetType == nil {
		return
	}

	// exclude generics, except specific types
	var isGeneric bool
	var generic *reflect.GenericType
	if generic, isGeneric = targetType.(*reflect.GenericType); isGeneric {
		if !includeGeneric(generic) {
			include = false
		}
	}

	// exclude built-in types
	if _, isBuiltIn := targetType.(*reflect.BuiltInType); isBuiltIn && !isGeneric {
		return
	}

	// exclude arrays
	if _, isArray := targetType.(*reflect.ArrayType); isArray {
		include = false
	}

	// exclude enums
	if isEnum(targetType) {
		include = false
	}

	// name := targetType.Name()
	// fmt.Println(name)

	typeName := targetType.FullName()
	if res.typesByName[typeName] != nil {
		return
	}

	res.typesByName[typeName] = targetType
	if include {
		res.typesFindOrder = append(res.typesFindOrder, targetType)
	}

	// generic arguments
	if genericType, ok := targetType.(*reflect.GenericType); ok {
		for _, argType := range genericType.ArgumentTypes() {
			res.innerResolve(argType)
		}
	}

	// base type
	if include {
		if baseType := targetType.Base(); baseType != nil {
			if !excludedBaseTypes[baseType.FullName()] {
				res.innerResolve(baseType)
			}
		}
	}

	// array value type
	if array, isArray := targetType.(*reflect.ArrayType); isArray {
		res.innerResolve(array.ValueType())
	}

	// methods
	if isServiceType {
		methods := targetType.GetMethods()
		for _, method := range methods {
			returnType := method.ReturnType()
			res.innerResolve(returnType)
			for _, param := range method.Parameters() {
				res.innerResolve(param.Type())
			}
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
	values := make([]reflect.Type, len(res.typesFindOrder))
	copy(values, res.typesFindOrder)

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
