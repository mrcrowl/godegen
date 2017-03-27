package reflect

import (
	"fmt"
	"godegen/cli"

	"github.com/bradfitz/slice"
	"github.com/ryanuber/go-glob"
)

type Assembly struct {
	metadata  *cli.Metadata
	typeCache *typeCache
	loader    *AssemblyLoader
}

func (asm *Assembly) LoadReferencedAssembly(assemblyName string) (*Assembly, error) {
	return asm.loader.Load(assemblyName)
}

func (asm *Assembly) Blah() []cli.IRow {
	table := asm.metadata.Tables.GetTable(cli.TableIdxAssembly)
	rows := table.Where(func(_ cli.IRow) bool { return true })
	return rows
}

func (asm *Assembly) GetType(name string) Type {
	t := asm.typeCache.get(name)
	if t == nil {
		if t = asm.loadType(name); t != nil {
			asm.typeCache.set(name, t)
		}
	}
	return t
}

func (asm *Assembly) GetTypesMatchingPattern(globPattern string, sorted bool) []Type {
	typeDefTable := asm.metadata.Tables.GetTable(cli.TableIdxTypeDef)
	typeMatchesPattern := func(row cli.IRow) bool {
		typeDef := row.(*cli.TypeDefRow)
		return glob.Glob(globPattern, typeDef.FullName())
	}
	typeRows := typeDefTable.Where(typeMatchesPattern)
	var matchingTypes []Type

	for _, row := range typeRows {
		typeDefRow := row.(*cli.TypeDefRow)
		matchedType := asm.typeCache.get(typeDefRow.FullName())
		if matchedType == nil {
			extendsType := asm.getTypeByIndex(typeDefRow.ExtendsIndex)
			matchedType = newTypeFromDef(typeDefRow, extendsType, asm)
		}
		matchingTypes = append(matchingTypes, matchedType)
	}

	// sort results
	if sorted {
		slice.Sort(matchingTypes, func(i, j int) bool {
			typeI := matchingTypes[i]
			typeJ := matchingTypes[j]
			if typeI.FullName() < typeJ.FullName() {
				return true
			}
			return false
		})
	}
	return matchingTypes
}

func (asm *Assembly) Test() {
	asm.metadata.Tables.GetTable(cli.TableIdxMethodSemantics).ForEach(func(row cli.IRow) {
		semRow := row.(*cli.MethodSemanticsRow)
		var text = "Prop "
		if semRow.Association.Type == 0 {
			return
		}
		fmt.Printf("%v\t%v\tKind\t%v\tMeth\t%v\n", text, semRow.Association.Row, semRow.Semantics, semRow.MethodRowNumber)
	})
}

func (asm *Assembly) GetTypeRowNumber(name string) uint32 {
	if rowNumber := asm.metadata.Tables.GetTable(cli.TableIdxTypeDef).RowNumberWhere(func(row cli.IRow) bool {
		return row.(*cli.TypeDefRow).FullName() == name
	}); rowNumber > 0 {
		return rowNumber
	}

	if rowNumber := asm.metadata.Tables.GetTable(cli.TableIdxTypeRef).RowNumberWhere(func(row cli.IRow) bool {
		return row.(*cli.TypeRefRow).FullName() == name
	}); rowNumber > 0 {
		return rowNumber
	}

	return 0
}

func (asm *Assembly) getTypeByIndex(index cli.TypeDefOrRefIndex) Type {
	var typ Type

	switch index.Type {
	case cli.TDORTypeDef:
		table := asm.metadata.Tables.GetTable(cli.TableIdxTypeDef)
		irow := table.GetRow(index.Row)
		if irow != nil {
			tdrow := irow.(*cli.TypeDefRow)
			fullName := tdrow.FullName()
			if typ = asm.typeCache.get(fullName); typ == nil {
				extendsType := asm.getTypeByIndex(tdrow.ExtendsIndex)
				typ = newTypeFromDef(tdrow, extendsType, asm)
				if typ != nil {
					asm.typeCache.set(typ.FullName(), typ)
				}
			}
		}

	case cli.TDORTypeRef:
		table := asm.metadata.Tables.GetTable(cli.TableIdxTypeRef)
		trrow := table.GetRow(index.Row).(*cli.TypeRefRow)
		fullName := trrow.FullName()
		if typ = asm.typeCache.get(fullName); typ == nil {
			typ = loadTypeFromRef(trrow, asm)
			if typ != nil {
				asm.typeCache.set(typ.FullName(), typ)
			}
		}

	case cli.TDORTypeSpec:
		// TODO:
		typ = nil
		// table := asm.metadata.Tables.GetTable(cli.TableIdxTypeSpec)
		// trrow := table.GetRow(index.Row).(*cli.TypeSpRow)
		// typ = newTypeFromRef(trrow, asm)
	}

	return typ
}

func (asm *Assembly) loadType(name string) Type {
	// type def
	typeDefTable := asm.metadata.Tables.GetTable(cli.TableIdxTypeDef)
	typeWithNameFn := func(row cli.IRow) bool {
		typeDef := row.(*cli.TypeDefRow)
		return typeDef.FullName() == name
	}
	if row := typeDefTable.First(typeWithNameFn); row != nil {
		typeDefRow := row.(*cli.TypeDefRow)
		extendsType := asm.getTypeByIndex(typeDefRow.ExtendsIndex)
		return newTypeFromDef(typeDefRow, extendsType, asm)
	}

	// type ref
	typeRefTable := asm.metadata.Tables.GetTable(cli.TableIdxTypeRef)
	typeRefWithNameFn := func(row cli.IRow) bool {
		typeRef := row.(*cli.TypeRefRow)
		return typeRef.FullName() == name
	}
	if row := typeRefTable.First(typeRefWithNameFn); row != nil {
		typeRefRow := row.(*cli.TypeRefRow)
		return loadTypeFromRef(typeRefRow, asm)
	}

	return nil
}
