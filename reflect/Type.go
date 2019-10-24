package reflect

import (
	"educationperfect.com/godegen/cli"
)

var importedNamespaces = map[string]bool{
	"System":                     true,
	"System.Threading.Tasks":     true,
	"System.Collections.Generic": true,
}

type Type interface {
	Name() string
	Namespace() string
	FullName() string
	Base() Type
	GetMethods() []*Method
	GetFields() []*Field
	GetFieldsWithOptions(includeNonPublic bool, includeInstance bool, includeStatic bool) []*Field
	GetProperties() []*Property
	RowNumber() uint32
}

func loadTypeFromRef(typeRefRow *cli.TypeRefRow, asm *Assembly) Type {
	tableSet := asm.metadata.Tables
	scopeIndex := typeRefRow.ScopeIndex

	switch typeRefRow.ScopeIndex.Type {
	case cli.RSModuleRef:
		// refRow := set.GetTable(TableIdxModuleRef).GetRow(row.ScopeIndex.Index).(*ModuleRefRow)
		// fmt.Println(refRow)
	case cli.RSAssemblyRef:
		assemblyRefRow := tableSet.GetTable(cli.TableIdxAssemblyRef).GetRow(scopeIndex.Index).(*cli.AssemblyRefRow)
		referenceName := assemblyRefRow.Name
		if referencedAssembly, error := asm.LoadReferencedAssembly(referenceName); error == nil {
			return referencedAssembly.GetType(typeRefRow.FullName())
		}

		// not loaded --> make a built-in type
		return newBuiltInTypeWithoutAlias(typeRefRow.FullName())

	default:
		panic("Not supported")
	}

	return nil
}
