package reflect

import "godegen/cli"

type Assembly struct {
	metadata  *cli.Metadata
	typeCache *typeCache
}

func LoadAssemblyFile(filepath string) *Assembly {
	assemblyPEFile := cli.OpenAssemblyPEFile(filepath)
	return &Assembly{assemblyPEFile.Metadata, newTypeCache()}
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

func (asm *Assembly) GetTypeByRow(row uint32) Type {
	table := asm.metadata.Tables.GetTable(cli.TableIdxTypeDef)
	row2 := table.GetRow(row).(*cli.TypeDefRow)
	return newTypeFromDef(row2, asm)
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
		return newTypeFromDef(typeDefRow, asm)
	}

	// type ref
	typeRefTable := asm.metadata.Tables.GetTable(cli.TableIdxTypeRef)
	typeRefWithNameFn := func(row cli.IRow) bool {
		typeRef := row.(*cli.TypeRefRow)
		return typeRef.FullName() == name
	}
	if row := typeRefTable.First(typeRefWithNameFn); row != nil {
		typeRefRow := row.(*cli.TypeRefRow)
		return newTypeFromRef(typeRefRow, asm)
	}

	return nil
}
