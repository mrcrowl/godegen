package main

import (
	"fmt"
	"godegen/description"
)

func main() {
	describer := description.NewServiceDescriber(`C:\WF\LP\server\EBS_Deployment\bin`, `Classes`)
	descr, _ := describer.Describe("nz.co.LanguagePerfect.Services.Portals.ControlPanel.LanguageDataPortal")
	fmt.Print(descr.JSON())
}
