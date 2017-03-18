package main

import (
	"godegen/description"
)

func main() {
	describer := description.NewServiceDescriber(`C:\WF\LP\server\EBS_Deployment\bin`, `Classes`)
	describer.Describe("nz.co.LanguagePerfect.Services.Portals.ControlPanel.LanguageDataPortal")
}
