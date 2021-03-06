package cli

import (
	"debug/pe"
)

// AssemblyPEFile is the logical representation of a CLI PE file
type AssemblyPEFile struct {
	*Metadata
}

// OpenAssemblyPEFile opens the specified filepath as an AssemblyPEFile
func OpenAssemblyPEFile(filepath string) (*AssemblyPEFile, error) {
	file, err := pe.Open(filepath)
	if err != nil {
		return nil, err
	}
	//is32Bit := (0x2000 & file.Characteristics) > 0
	optionalHeader, _ := file.OptionalHeader.(*pe.OptionalHeader32)
	cliDD := optionalHeader.DataDirectory[14]
	rawTextSection := file.Sections[0]

	textSection := newTextSection(rawTextSection)
	cliHeader := newHeader(textSection, cliDD)
	metadata := newMetadata(textSection, cliHeader.Metadata)
	return &AssemblyPEFile{metadata}, nil
}
