package codegen

import (
	"encoding/json"
	"errors"
	"fmt"
	"godegen/reflect"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type GeneratorConfig struct {
	Assembly                string                            `json:"assembly"`
	ServicePattern          []string                          `json:"servicePattern"`
	OutputPath              string                            `json:"outputPath"`
	TemplatesPath           string                            `json:"templatesPath"`
	FileExtension           string                            `json:"fileExtension"`
	TypeMap                 map[string]string                 `json:"typeMap"`
	NamespaceMap            map[string]string                 `json:"namespaceMap"`
	CollectionFormats       *GeneratorConfigCollectionFormats `json:"collectionFormats"`
	KeepServicesInNamespace bool                              `json:"keepServicesInNamespace"`
}

type GeneratorConfigCollectionFormats struct {
	System  string `json:"system"`
	Default string `json:"default"`
}

func LoadConfig(configFilename string) (*GeneratorConfig, error) {
	if filepath.Ext(configFilename) == "" {
		configFilename += ".json"
	}

	configBytes, e := ioutil.ReadFile(configFilename)
	if e != nil {
		return nil, errors.New("Invalid config file: " + configFilename)
	}

	var config *GeneratorConfig
	err := json.Unmarshal(configBytes, &config)
	if err != nil {
		return nil, errors.New("Invalid JSON in config file: " + err.Error())
	}
	return config, nil
}

func (config *GeneratorConfig) createTypeMapper() TypeMapperFn {
	var typeMap = config.TypeMap
	var namespaceMapperFn = config.createNamespaceMapper()
	var typeMapperFn TypeMapperFn
	var defaultFormat = config.CollectionFormats.Default
	var systemFormat = config.CollectionFormats.System
	var defaultFormatContainsVariable = strings.Contains(defaultFormat, "%")
	var systemFormatContainsVariable = strings.Contains(systemFormat, "%")

	typeMapperFn = func(typ reflect.Type) string {
		fullname := typ.FullName()

		if mappedName, found := typeMap[fullname]; found {
			return mappedName
		}

		cleanedFullName := namespaceMapperFn(fullname)
		if elemType, isCollection := isCollectionType(typ); isCollection {
			mappedTypeName := typeMapperFn(elemType)
			if isBuiltIn(elemType) {
				if defaultFormatContainsVariable {
					return fmt.Sprintf(systemFormat, mappedTypeName)
				}
				return systemFormat
			}
			if systemFormatContainsVariable {
				return fmt.Sprintf(defaultFormat, mappedTypeName)
			}
			return defaultFormat
		}

		return cleanedFullName
	}

	return typeMapperFn
}

func (config *GeneratorConfig) createNamespaceMapper() NamespaceMapperFn {
	namespaceMap := config.NamespaceMap
	return func(namespace string) string {
		for prefix, replacement := range namespaceMap {
			if strings.HasPrefix(namespace, prefix) {
				return replacement + namespace[len(prefix):]
			}
		}

		return namespace
	}
}

func isBuiltIn(typ reflect.Type) bool {
	_, isBuiltIn := typ.(*reflect.BuiltInType)
	return isBuiltIn
}

func isCollectionType(typ reflect.Type) (reflect.Type, bool) {
	if array, isArray := typ.(*reflect.ArrayType); isArray {
		return array.ValueType(), true
	}

	if generic, isGeneric := typ.(*reflect.GenericType); isGeneric {
		if generic.TypeBase.FullName() == "System.Collections.Generic.List`1" {
			return generic.ArgumentTypes()[0], true
		}
	}

	return nil, false
}
