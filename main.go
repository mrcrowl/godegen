package main

import (
	"bytes"
	"fmt"
	"godegen/cli"
	"godegen/reflect"
)

func main() {
	assemblyFile := reflect.LoadAssemblyFile(`C:\WF\LP\server\EBS_Deployment\bin\Classes.dll`)

	utp := assemblyFile.GetType("nz.co.LanguagePerfect.Services.UserTasks.Managers.UserTaskManager").(*reflect.TypeDef)

	methods := utp.GetMethods()
	firstMethod := methods[0]
	fmt.Println(firstMethod.Blob)

	type1 := assemblyFile.GetTypeByRow(529)  //529)
	type2 := assemblyFile.GetTypeByRow(1135) //1135)
	fmt.Println(type1)
	fmt.Println(type2)

	row1 := assemblyFile.GetTypeRowNumber("nz.co.LanguagePerfect.Services.Sessions.BusinessObjects.LPSession")
	fmt.Println(row1)

	sr := cli.NewShapeReader(bytes.NewReader([]byte{0xC0, 00, 0x40, 00}))
	test := sr.ReadCompressedUInt()
	fmt.Printf("0x%x", test)

	// fmt.Println("--TYPEDEF--")
	// for _, row := range metadata.Tables.GetRows(cli.TableIdxTypeDef) {
	// 	fmt.Println(row)
	// }

	// fmt.Println()
	// fmt.Println("--FIELD--")

	// for _, row := range metadata.Tables.GetRows(cli.TableIdxField) {
	// 	fmt.Println(row)
	// }

	// fmt.Println()
	// fmt.Println("--METHODDEF--")

	// for _, row := range assemblyFile.Tables.GetRows(cli.TableIdxMethodDef) {
	// 	fmt.Println(row)
	// }
}
