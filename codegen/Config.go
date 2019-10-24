package codegen

import (
	"encoding/json"
	"errors"
	"fmt"
	"educationperfect.com/godegen/reflect"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// GeneratorConfig is the internal representation of the JSON configuration file
type GeneratorConfig struct {
	Assembly                string            `json:"assembly"`
	ServicePattern          []string          `json:"servicePattern"`
	OutputPath              string            `json:"outputPath"`
	TemplatesPath           string            `json:"templatesPath"`
	DataTypePathSubfolder   string            `json:"dataTypePathSubfolder"`
	FileExtension           string            `json:"fileExtension"`
	TypeMap                 map[string]string `json:"typeMap"`
	NamespaceMap            map[string]string `json:"namespaceMap"`
	CollectionFormats       map[string]string `json:"collectionFormats"`
	KeepServicesInNamespace bool              `json:"keepServicesInNamespace"`
	ServiceRelocationMap    map[string]string `json:"serviceRelocationMap"`
}

// GeneratorConfigCollectionFormats are the collection formats used for arrays/lists
type GeneratorConfigCollectionFormats struct {
	System  string `json:"system"`
	Default string `json:"default"`
}

// LoadConfig loads a JSON configuration file
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

func (config *GeneratorConfig) getEnumType() string {
	if mappedName, found := config.TypeMap["System.Enum"]; found {
		return mappedName
	}

	return "string"
}

func (config *GeneratorConfig) createTypeMapper() typeMapperFunc {
	var typeMap = config.TypeMap
	var namespaceMapperFn = config.createNamespaceMapper()
	var typeMapperFn typeMapperFunc
	var defaultFormat = config.CollectionFormats["default"]
	var systemFormat = config.CollectionFormats["system"]
	var defaultFormatContainsVariable = strings.Contains(defaultFormat, "%")
	var systemFormatContainsVariable = strings.Contains(systemFormat, "%")

	typeMapperFn = func(typ reflect.Type, nameOnly bool) string {
		fullname := typ.FullName()

		if mappedName, found := typeMap[fullname]; found {
			return mappedName
		}

		if elemType, isCollection := isCollectionType(typ); isCollection {
			mappedTypeName := typeMapperFn(elemType, nameOnly)
			if isBuiltIn(elemType) {
				if systemFormatContainsVariable {
					return fmt.Sprintf(systemFormat, mappedTypeName)
				}
				return systemFormat
			}
			if defaultFormatContainsVariable {
				return fmt.Sprintf(defaultFormat, mappedTypeName)
			}
			return defaultFormat
		}

		// special handling for enum
		if isEnum(typ) {
			return config.getEnumType()
		}

		if nameOnly {
			return typ.Name()
		}
		mappedFullName := namespaceMapperFn(fullname)
		return mappedFullName
	}

	return typeMapperFn
}

func (config *GeneratorConfig) createNamespaceMapper() namespaceMapperFunc {
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

func isCollection(typ reflect.Type) bool {
	if _, isArray := typ.(*reflect.ArrayType); isArray {
		return true
	}

	if generic, isGeneric := typ.(*reflect.GenericType); isGeneric {
		if generic.TypeBase.FullName() == "System.Collections.Generic.List`1" {
			return true
		}
	}

	return false
}

func isEnum(typ reflect.Type) bool {
	// special handling for enum
	base := typ.Base()
	if base != nil && base.FullName() == "System.Enum" {
		return true
	}
	return false
}

func isGeneric(typ reflect.Type) bool {
	if _, isGeneric := typ.(*reflect.GenericType); isGeneric {
		return true
	}

	return false
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
