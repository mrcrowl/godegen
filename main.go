package main

import (
	"fmt"
	"godegen/reflect"
)

func main() {
	assemblyFile := reflect.LoadAssemblyFile(`C:\WF\LP\server\EBS_Deployment\bin\Classes.dll`)

	utp := assemblyFile.GetType("nz.co.LanguagePerfect.Services.UserTasks.Managers.UserTaskManager").(*reflect.TypeDef)

	methods := utp.GetMethods()
	firstMethod := methods[0]
	fmt.Println(firstMethod)

	row1 := assemblyFile.GetTypeRowNumber("nz.co.LanguagePerfect.Services.Sessions.BusinessObjects.LPSession")
	fmt.Println(row1)

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
