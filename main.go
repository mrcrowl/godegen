package main

import (
	"debug/pe"
	"fmt"
	"godegen/cli"
)

func main() {
	file, _ := pe.Open(`C:\WF\LP\server\EBS_Deployment\bin\Classes.dll`)
	//is32Bit := (0x2000 & file.Characteristics) > 0
	optionalHeader, _ := file.OptionalHeader.(*pe.OptionalHeader32)
	cliDD := optionalHeader.DataDirectory[14]
	fmt.Printf("%v", cliDD.VirtualAddress)
	rawTextSection := file.Sections[0]

	textSection := cli.NewTextSection(rawTextSection)
	cliHeader := cli.NewHeader(textSection, cliDD)
	metadata := cli.NewMetadata(textSection, cliHeader.Metadata)
	// magic := binary.LittleEndian.Uint32(metaDataBuffer.Bytes())
	// hexString := strconv.FormatUint(uint64(magic), 16)
	// fmt.Printf(hexString)

	// fmt.Println("--TYPEDEF--")
	// for _, row := range metadata.Tables.GetRows(cli.TableIdxTypeDef) {
	// 	fmt.Println(row)
	// }

	// fmt.Println()
	// fmt.Println("--FIELD--")

	// for _, row := range metadata.Tables.GetRows(cli.TableIdxField) {
	// 	fmt.Println(row)
	// }

	fmt.Println()
	fmt.Println("--METHODDEF--")

	for _, row := range metadata.Tables.GetRows(cli.TableIdxMethodDef) {
		fmt.Println(row)
	}
}
