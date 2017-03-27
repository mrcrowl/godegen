package codegen

import (
	"errors"
	"godegen/reflect"
	"path/filepath"
	"strings"

	"github.com/bradfitz/slice"
)

type TypeMapperFn func(reflect.Type) string
type NamespaceMapperFn func(string) string

type ServiceDescriber struct {
	assemblyFile    *reflect.Assembly
	typeMapper      TypeMapperFn
	namespaceMapper NamespaceMapperFn
}

func (descr *ServiceDescriber) GetTypesMatchingPattern(globPattern string) []reflect.Type {
	return descr.assemblyFile.GetTypesMatchingPattern(globPattern, true)
}

func (descr *ServiceDescriber) GetType(typeName string) []reflect.Type {
	return []reflect.Type{descr.assemblyFile.GetType(typeName)}
}

func NewServiceDescriber(config *GeneratorConfig) (*ServiceDescriber, error) {
	assemblyPath, assemblyName := filepath.Split(config.Assembly)

	return NewServiceDescriberManual(
		assemblyPath,
		assemblyName,
		config.createTypeMapper(),
		config.createNamespaceMapper(),
	)
}

func NewServiceDescriberManual(
	assemblyPath string,
	assemblyName string,
	typeMapper func(reflect.Type) string,
	namespaceMapper func(string) string,
) (*ServiceDescriber, error) {
	loader := reflect.NewAssemblyLoader(assemblyPath)
	assemblyFile, err := loader.Load(assemblyName)
	if err == nil {
		return &ServiceDescriber{assemblyFile, typeMapper, namespaceMapper}, nil
	}

	return nil, errors.New("Can't load assembly '" + assemblyName + "' in: " + assemblyPath)
}

func (res *ServiceDescriber) Describe(serviceTypeName string) (*ServiceDescription, error) {
	if serviceType := res.assemblyFile.GetType(serviceTypeName); serviceType != nil {
		return res.DescribeType(serviceType)
	}

	return nil, errors.New("Can't find type: " + serviceTypeName)
}

func (res *ServiceDescriber) DescribeType(serviceType reflect.Type) (*ServiceDescription, error) {
	resolvedTypes := ResolveServiceDependencyTypes(res.assemblyFile, serviceType)
	description := res.createDescriptionOfTypes(resolvedTypes, serviceType)
	return description, nil
}

func (res *ServiceDescriber) createDescriptionOfTypes(types []reflect.Type, serviceType reflect.Type) *ServiceDescription {
	rootNamespaces, namespaceMap, referencedNamespaces := res.buildNamespaceTree(types, serviceType)

	// add service
	serviceTypeNamespace := res.mapNamespace(serviceType.Namespace())
	serviceNamespace := namespaceMap[serviceTypeNamespace]
	service := res.createService(serviceType, referencedNamespaces)
	serviceNamespace.addService(service)

	// fmt.Println(rootNamespaces)
	return &ServiceDescription{
		Namespaces: rootNamespaces,
	}
}

func (res *ServiceDescriber) createService(serviceType reflect.Type, referencedNamespaces []string) *Service {
	methods := res.collectTypeMethods(serviceType)
	serviceIdentifier := serviceType.FullName()

	slice.Sort(referencedNamespaces, func(i, j int) bool {
		nsi := referencedNamespaces[i]
		nsj := referencedNamespaces[j]
		return nsi < nsj
	})

	return &Service{
		*res.createDataType(serviceType),
		serviceIdentifier,
		methods,
		referencedNamespaces,
	}
}

func (res *ServiceDescriber) buildNamespaceTree(types []reflect.Type, serviceType reflect.Type) ([]*Namespace, map[string]*Namespace, []string) {
	var namespaceSeen = map[string]*Namespace{}
	var distinctNamespaces []*Namespace
	var rootNamespaces []*Namespace
	var namespacesWithTypes = make(map[string]bool)

	for _, typ := range append(types, serviceType) {
		var childNamespace *Namespace
		var namespace *Namespace
		var found bool
		var nsName = res.mapNamespace(typ.Namespace())
		var nsOriginal = nsName
		for {
			if namespace, found = namespaceSeen[nsName]; !found {
				namespace = newNamespace(nsName)
			}

			if childNamespace != nil {
				namespace.addChild(childNamespace)
			}

			if found {
				break
			}

			namespaceSeen[nsName] = namespace
			distinctNamespaces = append(distinctNamespaces, namespace)
			if namespace.isRoot() {
				rootNamespaces = append(rootNamespaces, namespace)
				break
			}

			nsName = getParentNamespace(nsName)
			childNamespace = namespace
		}

		namespace = namespaceSeen[nsOriginal]

		// add type to namespace
		if typ != serviceType {
			dataType := res.createDataType(typ)
			namespace.DataTypes = append(namespace.DataTypes, dataType)
			namespacesWithTypes[namespace.qualifiedName] = true
		}
	}

	referencedNamespaces := make([]string, 0, len(namespacesWithTypes))

	for namespace, _ := range namespacesWithTypes {
		referencedNamespaces = append(referencedNamespaces, namespace)
	}

	return rootNamespaces, namespaceSeen, referencedNamespaces
}

func (res *ServiceDescriber) createDataType(typ reflect.Type) *DataType {
	var fields = res.collectTypeFields(typ)
	var consts = res.collectConsts(typ)
	var baseType = typ.Base()
	var base *DataTypeReference

	if !excludedBaseTypes[baseType.FullName()] {
		base = res.createDataTypeReference(baseType)
	}

	return &DataType{
		DataTypeReference{
			Name:      typ.Name(),
			Namespace: res.mapNamespace(typ.Namespace()),
			// QualifiedName: res.mapNamespace(typ.FullName()),
		},
		base, // TODO: get base type
		fields,
		consts,
	}
}

func (res *ServiceDescriber) createDataTypeReference(typ reflect.Type) *DataTypeReference {
	if typ == nil {
		return nil
	}

	return &DataTypeReference{
		Name:      typ.Name(),
		Namespace: res.mapNamespace(typ.Namespace()),
	}
}

func (res *ServiceDescriber) collectTypeFields(typ reflect.Type) []*Field {
	var fields []*Field

	for _, property := range typ.GetProperties() {
		fields = append(fields, res.createFieldFromProperty(property))
	}

	for _, field := range typ.GetFields() {
		fields = append(fields, res.createFieldFromField(field))
	}

	return fields
}

func (res *ServiceDescriber) collectConsts(typ reflect.Type) []*Const {
	var consts []*Const

	for _, constant := range typ.GetFieldsWithOptions(false, false, true) {
		if constant.Value() != nil {
			consts = append(consts, res.createConst(constant))
		}
	}

	return consts
}

func (res *ServiceDescriber) mapType(typ reflect.Type) string {
	return res.typeMapper(typ)
}

func (res *ServiceDescriber) mapElementType(typ reflect.Type) string {
	if elementType, is := isCollectionType(typ); is {
		return res.mapType(elementType)
	}

	return ""
}

func (res *ServiceDescriber) mapNamespace(namespace string) string {
	return res.namespaceMapper(namespace)
}

func (res *ServiceDescriber) collapseReturnType(returnType reflect.Type) reflect.Type {
	if genericType, isGeneric := returnType.(*reflect.GenericType); isGeneric {
		if genericType.Namespace() == "System.Threading.Tasks" &&
			genericType.LexicalName() == "Task" {
			return genericType.ArgumentTypes()[0]
		}
	}

	return returnType
}

func (res *ServiceDescriber) collectTypeMethods(typ reflect.Type) []*Method {
	var methods []*Method

	for _, method := range typ.GetMethods() {
		args := res.collectMethodArgs(method)
		returnType := res.collapseReturnType(method.ReturnType())

		meth := &Method{
			Name:     method.Name(),
			Type:     res.mapType(returnType),
			Args:     args,
			nameSort: strings.ToLower(method.Name()),
		}
		methods = append(methods, meth)
	}

	slice.Sort(methods, func(i, j int) bool {
		methodI := methods[i]
		methodJ := methods[j]
		if methodI.nameSort < methodJ.nameSort {
			return true
		}
		return false
	})

	return methods
}

func (res *ServiceDescriber) collectMethodArgs(method *reflect.Method) []*Arg {
	var args []*Arg

	for _, param := range method.Parameters() {
		arg := &Arg{
			Name: param.Name(),
			Type: res.mapType(param.Type()),
		}
		args = append(args, arg)
	}

	return args
}

func (res *ServiceDescriber) createConst(field *reflect.Field) *Const {
	return &Const{
		Name:  field.Name(),
		Type:  res.mapType(field.Type()),
		Value: field.Value(),
	}
}

func (res *ServiceDescriber) createFieldFromField(field *reflect.Field) *Field {
	return &Field{
		Name:        field.Name(),
		Type:        res.mapType(field.Type()),
		ElementType: res.mapElementType(field.Type()),
	}
}

func (res *ServiceDescriber) createFieldFromProperty(property *reflect.Property) *Field {
	return &Field{
		Name:        property.Name(),
		Type:        res.mapType(property.Type()),
		ElementType: res.mapElementType(property.Type()),
	}
}

func newNamespace(namespace string) *Namespace {
	name := getLastNamespaceSegment(namespace)
	return &Namespace{
		Name:          name,
		qualifiedName: namespace,
	}
}

func getLastNamespaceSegment(namespace string) string {
	dot := strings.LastIndex(namespace, ".")
	if dot >= 0 {
		return namespace[dot+1:]
	}
	return namespace
}

func getParentNamespace(namespace string) string {
	dot := strings.LastIndex(namespace, ".")
	if dot < 0 {
		return ""
	}
	return namespace[:dot]
}

// func printTypes(resolvedTypes []reflect.Type) {
// 	for _, t := range resolvedTypes {
// 		fmt.Println(t.FullName() + ":")
// 		for _, field := range t.GetFields() {
// 			fieldType := field.Type()
// 			if collectionType, isCollection := isCollectionType(fieldType); isCollection {
// 				fmt.Println("\t" + field.Name() + ": " + collectionType.FullName() + "[]")
// 			} else {
// 				fmt.Println("\t" + field.Name() + ": " + field.Type().FullName())
// 			}
// 		}
// 	}
// }
