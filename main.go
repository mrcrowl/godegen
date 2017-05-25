package main

import (
	"flag"
	"fmt"
	"godegen/codegen"
	"os"
	"strings"
)

var maxServicePatternLen = 50

func usage() {
	fmt.Println("Usage: codegen.exe  [-service=glob] <configFile>")
}

func main() {
	serviceName := flag.String("service", "", "")
	flag.Parse()

	configName := flag.Arg(0)
	if configName == "" {
		usage()
		os.Exit(0)
	}

	config, err := codegen.LoadConfig(configName)
	if err != nil {
		fmt.Println(err)
		os.Exit(4)
	}

	if (*serviceName) != "" {
		config.ServicePattern = []string{*serviceName}
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
		numChanges := 0
		serviceTypes := describer.GetTypesMatchingPattern(servicePattern)
		for _, serviceType := range serviceTypes {
			serviceName := getServiceName(serviceType.FullName())
			if !strings.Contains(serviceName, "PortalsAsync") {
				fmt.Println("ERROR: Not an async portal!")
				continue
			}
			fmt.Printf("\n%s ", serviceName)
			descr, _ := describer.DescribeType(serviceType)

			if numChanges, err = gen.OutputServiceDescription(descr); err != nil {
				fmt.Println(err)
				os.Exit(2)
			}

			if numChanges == 0 {
				paddingLen := maxServicePatternLen - len(serviceName)
				fmt.Print(leftPad(" no change", "-", paddingLen))
			} else {
				fmt.Println()
			}
		}
	}
}

func leftPad(s string, padStr string, pLen int) string {
	if pLen < (len(padStr) + 3) {
		pLen = len(padStr) + 3
	}
	return strings.Repeat(padStr, pLen) + s
}

func getServiceName(serviceTypeName string) string {
	parts := strings.Split(serviceTypeName, ".")
	if len(parts) >= 3 {
		return strings.Join(parts[len(parts)-3:], ".")
	}
	return serviceTypeName
}
