package main

import (
	"fmt"
	"godegen/codegen"
	"os"
)

func usage() {
	fmt.Println("Usage: godegen.exe -f  <configFile>")
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(3)
	}
	configName := os.Args[1]
	config, err := codegen.LoadConfig(configName)
	if err != nil {
		fmt.Println(err)
		os.Exit(4)
	}

	describer, err := codegen.NewServiceDescriber(config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	gen, err := codegen.NewGenerator(config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, servicePattern := range config.ServicePattern {
		serviceTypes := describer.GetTypesMatchingPattern(servicePattern)
		for _, serviceType := range serviceTypes {
			fmt.Println(serviceType.FullName())
			descr, _ := describer.DescribeType(serviceType)

			if err = gen.OutputServiceDescription(descr); err != nil {
				fmt.Println(err)
				os.Exit(2)
			}
		}
	}
}
