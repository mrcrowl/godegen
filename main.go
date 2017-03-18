package main

import (
	"fmt"
	"godegen/description"
	"godegen/reflect"
)

func main() {
	loader := reflect.NewAssemblyLoader(`C:\WF\LP\server\EBS_Deployment\bin`)
	assemblyFile, _ := loader.Load("Classes")
	// assemblyFile.Test()
	//assemblyFile, _ := loader.Load("LiteCASClient")

	// type1 := assemblyFile.GetType("nz.co.LanguagePerfect.Services.Classes.BusinessObjects.ClassDescription")
	str := description.NewServiceTypesResolver(assemblyFile)
	// ssoRego := assemblyFile.GetType("nz.co.LanguagePerfect.Services.LPLogin.Managers.SSORegistrationManager")
	ssoRego := assemblyFile.GetType("nz.co.LanguagePerfect.Services.Portals.ControlPanel.LanguageDataPortal")
	resolvedTypes := str.Resolve(ssoRego)
	for _, t := range resolvedTypes {
		fmt.Println(t.FullName())
	}

	// utp := assemblyFile.GetType("nz.co.LanguagePerfect.Services.Sessions.BusinessObjects.LPSession").(*reflect.TypeDef)

	// methods := utp.GetMethods()
	// for _, method := range methods {
	// 	fmt.Println(method)
	// }

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
