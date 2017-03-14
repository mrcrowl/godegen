package cli

import "encoding/binary"
import "debug/pe"

// Header =
type Header struct {
	Size                    uint32
	MajorRuntime            uint16
	MinorRuntime            uint16
	Metadata                RVA
	Flags                   uint32
	EntryPointToken         uint32
	Resources               RVA
	StrongNameSig           RVA
	CodeManageTable         RVA
	VtableFixups            RVA
	ExportAddressTableJumps RVA
	ManagedNativeHeader     RVA
}

func newHeader(textSection *TextSection, cliDD pe.DataDirectory) *Header {
	cliHeaderBuffer := textSection.GetRange(cliDD.VirtualAddress, cliDD.Size)
	cliHeader := new(Header)
	binary.Read(cliHeaderBuffer, binary.LittleEndian, cliHeader)
	return cliHeader
}
