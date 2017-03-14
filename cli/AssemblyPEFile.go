package cli

import (
	"debug/pe"
	"fmt"
)

// AssemblyPEFile is the logical representation of a CLI PE file
type AssemblyPEFile struct {
	*Metadata
}

// OpenAssemblyPEFile opens the specified filepath as an AssemblyPEFile
func OpenAssemblyPEFile(filepath string) *AssemblyPEFile {
	file, _ := pe.Open(filepath)
	//is32Bit := (0x2000 & file.Characteristics) > 0
	optionalHeader, _ := file.OptionalHeader.(*pe.OptionalHeader32)
	cliDD := optionalHeader.DataDirectory[14]
	fmt.Printf("%v", cliDD.VirtualAddress)
	rawTextSection := file.Sections[0]

	textSection := newTextSection(rawTextSection)
	cliHeader := newHeader(textSection, cliDD)
	metadata := newMetadata(textSection, cliHeader.Metadata)
	return &AssemblyPEFile{metadata}
}
