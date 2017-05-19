package codegen

import (
	"errors"
	"godegen/reflect"
	"path/filepath"
	"strings"

	"strconv"

	"github.com/bradfitz/slice"
)

type typeMapperFunc func(reflect.Type, bool) string
type namespaceMapperFunc func(string) string

type ServiceDescriber struct {
	assemblyFile    *reflect.Assembly
	typeMapper      typeMapperFunc
	namespaceMapper namespaceMapperFunc
}

func (res *ServiceDescriber) GetTypesMatchingPattern(globPattern string) []reflect.Type {
	return res.assemblyFile.GetTypesMatchingPattern(globPattern, true)
}

func (res *ServiceDescriber) GetType(typeName string) []reflect.Type {
	return []reflect.Type{res.assemblyFile.GetType(typeName)}
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
	typeMapper func(reflect.Type, bool) string,
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

	serviceDataType := res.createDataType(serviceType, true)
	if serviceDataType.aliasMap.nonEmpty() {
		for _, m := range methods {
			m.ApplyAliases(serviceDataType.aliasMap)
		}
	}

	return &Service{
		*serviceDataType,
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
			dataType := res.createDataType(typ, false)
			namespace.DataTypes = append(namespace.DataTypes, dataType)
			namespacesWithTypes[namespace.qualifiedName] = true
		}
	}

	referencedNamespaces := make([]string, 0, len(namespacesWithTypes))

	for namespace := range namespacesWithTypes {
		referencedNamespaces = append(referencedNamespaces, namespace)
	}

	return rootNamespaces, namespaceSeen, referencedNamespaces
}

func (res *ServiceDescriber) createDataType(typ reflect.Type, includeMethods bool) *DataType {
	var fields = res.collectTypeFields(typ)
	var consts = res.collectConsts(typ)
	var baseType = typ.Base()
	var base *RelativeDataTypeReference

	if !excludedBaseTypes[baseType.FullName()] {
		base = res.createRelativeDataTypeReference(baseType, typ)
	}

	referencedTypes := res.collectReferencedTypes(typ, includeMethods)
	aliasMap := res.createDuplicateNameAliasMap(referencedTypes)

	if aliasMap.nonEmpty() {
		for _, f := range fields {
			f.ApplyAliases(aliasMap)
		}
		for _, c := range consts {
			c.ApplyAliases(aliasMap)
		}
	}

	return &DataType{
		DataTypeReference{
			Name:      typ.Name(),
			Namespace: res.mapNamespace(typ.Namespace()),
			// QualifiedName: res.mapNamespace(typ.FullName()),
		},
		base,
		referencedTypes,
		fields,
		consts,
		aliasMap,
	}
}

func (res *ServiceDescriber) createRelativeDataTypeReference(typ reflect.Type, relativeTo reflect.Type) *RelativeDataTypeReference {
	if typ == nil {
		return nil
	}

	name := typ.Name()
	namespace := res.mapNamespace(typ.Namespace())

	dataTypeRef := &DataTypeReference{name, namespace}

	var relativePath string
	if isBuiltIn(typ) || isBuiltIn(relativeTo) || isGeneric(typ) || isGeneric(relativeTo) {
		relativePath = ""
	} else {
		fromPath := strings.Replace(namespace, ".", "/", -1)
		toNamespace := res.mapNamespace(relativeTo.Namespace())
		toPath := strings.Replace(toNamespace, ".", "/", -1)
		// println(fromPath)
		// println(toPath)
		relativePath = calculateRelativePath(fromPath, toPath)
		// println(relativePath)
	}
	return &RelativeDataTypeReference{
		DataTypeReference: *dataTypeRef,
		RelativePath:      relativePath,
	}
}

func calculateRelativePath(fromPath string, toPath string) string {
	if fromPath == toPath {
		return "."
	}

	if relativePath, err := filepath.Rel(toPath, fromPath); err == nil {
		slashedPath := filepath.ToSlash(relativePath)
		if strings.HasPrefix(slashedPath, ".") {
			return slashedPath
		}
		return "./" + slashedPath
	}
	return ""
}

func (res *ServiceDescriber) createDuplicateNameAliasMap(referencedTypes []*RelativeDataTypeReference) aliasMap { // map[fullName] --> alias
	// group by name
	var typesByName = map[string][]*RelativeDataTypeReference{}
	var duplicateNames = map[string]bool{}

	for _, ref := range referencedTypes {
		name := ref.Name
		existing := append(typesByName[name], ref)
		typesByName[name] = existing

		// keep track of duplicates
		if len(existing) > 1 {
			duplicateNames[name] = true
		}
	}

	var aliasesByFullname = aliasMap{}
	for duplicateName := range duplicateNames {
		references := typesByName[duplicateName]
		for i, ref := range references {
			fullName := ref.Namespace + "." + ref.Name
			alias := ref.Name + "_" + strconv.Itoa(i+1)
			aliasesByFullname[fullName] = alias
			ref.Alias = alias
		}
	}

	return aliasesByFullname
}

func (res *ServiceDescriber) collectReferencedTypes(sourceType reflect.Type, includeMethods bool) []*RelativeDataTypeReference {
	var types []reflect.Type
	var typesSeen = map[string]bool{}

	var collect func(reflect.Type)
	collect = func(typ reflect.Type) {
		if isBuiltIn(typ) {
			return
		}

		if isEnum(typ) {
			return
		}

		if generic, isGeneric := typ.(*reflect.GenericType); isGeneric {
			fullName := generic.TypeBase.FullName()

			if _, inWhiteList := genericsWhitelist[fullName]; inWhiteList {
				collect(generic.ArgumentTypes()[0])
			}
			return
		}

		if elementType, isCollection := isCollectionType(typ); isCollection {
			collect(elementType)
			return
		}

		fullname := typ.FullName()
		if _, seen := typesSeen[fullname]; seen {
			return
		}

		types = append(types, typ)
		typesSeen[fullname] = true
	}

	// base type
	var baseType = sourceType.Base()
	if !excludedBaseTypes[baseType.FullName()] {
		collect(baseType)
	}

	// method
	if includeMethods {
		for _, method := range sourceType.GetMethods() {
			collect(method.ReturnType())
			for _, param := range method.Parameters() {
				collect(param.Type())
			}
		}
	}

	// properties
	for _, property := range sourceType.GetProperties() {
		collect(property.Type())
	}

	// fields
	for _, field := range sourceType.GetFields() {
		collect(field.Type())
	}

	var dataTypeReferences []*RelativeDataTypeReference
	for _, typ := range types {
		ref := res.createRelativeDataTypeReference(typ, sourceType)
		dataTypeReferences = append(dataTypeReferences, ref)
	}

	return dataTypeReferences
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

func (res *ServiceDescriber) mapType(typ reflect.Type, nameOnly bool) string {
	mappedName := res.typeMapper(typ, nameOnly)
	return mappedName
}

func (res *ServiceDescriber) mapElementType(typ reflect.Type, nameOnly bool) string {
	if elementType, is := isCollectionType(typ); is {
		return res.mapType(elementType, nameOnly)
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
			Type:     res.mapType(returnType, false),
			TypeName: res.mapType(returnType, true),
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
			Name:     param.Name(),
			Type:     res.mapType(param.Type(), false),
			TypeName: res.mapType(param.Type(), true),
		}
		args = append(args, arg)
	}

	return args
}

func (res *ServiceDescriber) createConst(field *reflect.Field) *Const {
	return &Const{
		Name:     field.Name(),
		Type:     res.mapType(field.Type(), false),
		TypeName: res.mapType(field.Type(), true),
		Value:    field.Value(),
	}
}

func (res *ServiceDescriber) createFieldFromField(field *reflect.Field) *Field {
	fieldType := field.Type()

	return &Field{
		Name:            field.Name(),
		Type:            res.mapType(fieldType, false),
		TypeName:        res.mapType(fieldType, true),
		ElementType:     res.mapElementType(fieldType, false),
		ElementTypeName: res.mapElementType(fieldType, true),
	}
}

func (res *ServiceDescriber) createFieldFromProperty(property *reflect.Property) *Field {
	propertyType := property.Type()

	return &Field{
		Name:            property.Name(),
		Type:            res.mapType(propertyType, false),
		TypeName:        res.mapType(propertyType, true),
		ElementType:     res.mapElementType(propertyType, false),
		ElementTypeName: res.mapElementType(propertyType, true),
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
