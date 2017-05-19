package codegen

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	templateNameDataType = "datatype"
	templateNameService  = "service"
)

var additionalTemplateFns = template.FuncMap{
	"replace": func(s, old, new string) string {
		return strings.Replace(s, old, new, -1)
	},
	"typeName": func(fullyQualifiedTypeName string) string {
		lastDot := strings.LastIndexByte(fullyQualifiedTypeName, '.')
		if lastDot >= 0 {
			return fullyQualifiedTypeName[lastDot+1:]
		}
		return fullyQualifiedTypeName
	},
	"namespaceName": func(fullyQualifiedTypeName string) string {
		lastDot := strings.LastIndexByte(fullyQualifiedTypeName, '.')
		if lastDot >= 0 {
			return fullyQualifiedTypeName[:lastDot]
		}
		return fullyQualifiedTypeName
	},
}

// Generator is used to apply a set of templates to a ServiceDescription and output to files
type Generator struct {
	outputPath string
	config     *GeneratorConfig
	templates  *template.Template
}

// NewGenerator creates a new Generator
func NewGenerator(config *GeneratorConfig) (*Generator, error) {
	// validate output path
	outputPath := config.OutputPath
	if !isValidDirectory(outputPath) {
		return nil, errors.New("Invalid output path: " + outputPath)
	}

	// validate templates
	if !isValidDirectory(config.TemplatesPath) {
		return nil, errors.New("Invalid templates path: " + config.TemplatesPath)
	}

	templatePattern := filepath.Join(config.TemplatesPath, "*.gotmpl")
	templates, err := template.New("main").Funcs(additionalTemplateFns).ParseGlob(templatePattern)
	if err != nil {
		return nil, err
	}

	for _, templateName := range []string{templateNameDataType, templateNameService} {
		if templates.Lookup(templateName) == nil {
			return nil, errors.New("Could not find template '" + templateName + "' in: " + config.TemplatesPath)
		}
	}

	return &Generator{
		outputPath,
		config,
		templates,
	}, nil
}

// OutputServiceDescription writes a service description to file
func (gen *Generator) OutputServiceDescription(descr *ServiceDescription) (int, error) {
	totalChanges := 0
	for _, namespace := range descr.Namespaces {
		numChanges, err := gen.outputNamespace(gen.outputPath, namespace)
		totalChanges += numChanges
		if err != nil {
			return totalChanges, err
		}
	}

	return totalChanges, nil
}

func (gen *Generator) outputNamespace(outputPath string, namespace *Namespace) (int, error) {
	numChanges := 0
	namespacePath := joinNamespaceToOutputPath(outputPath, namespace.Name)
	os.Mkdir(namespacePath, os.ModePerm)

	// datatype
	for _, dataType := range namespace.DataTypes {
		filename, changed, err := gen.outputDataType(namespacePath, dataType)
		if err != nil {
			return 0, err
		}
		if changed {
			numChanges++
			fmt.Printf("\n - %s.%s", namespace.Name, filename)
		}
	}

	// service
	for _, service := range namespace.Services {
		var servicePath string
		if gen.config.KeepServicesInNamespace {
			servicePath = namespacePath
		} else {
			// hacky workaround for changing portals output path
			servicePath = gen.config.OutputPath
			for _, typ := range service.ReferencedTypes {
				typ.RelativePath = strings.Replace(typ.RelativePath, "../../", "./", -1)
			}
		}
		filename, changed, err := gen.outputService(servicePath, service)
		if err != nil {
			return 0, err
		}
		if changed {
			numChanges++
			fmt.Printf("\n - %s", filename)
		}
	}

	// sub-namespaces
	for _, namespace := range namespace.Namespaces {
		numSubChanges, err := gen.outputNamespace(namespacePath, namespace)
		numChanges += numSubChanges
		if err != nil {
			return 0, err
		}
	}

	return numChanges, nil
}

func (gen *Generator) outputDataType(outputPath string, dataType *DataType) (string, bool, error) {
	var dataTypeFilename = dataType.Name + gen.config.FileExtension
	var dataTypePath string
	if gen.config.DataTypePathSubfolder != "" {
		dataTypePath = filepath.Join(outputPath, gen.config.DataTypePathSubfolder)
		os.MkdirAll(dataTypePath, os.ModePerm)
	} else {
		dataTypePath = outputPath
	}
	dataTypeFilePath := filepath.Join(dataTypePath, dataTypeFilename)

	changed, err := gen.outputTemplateIfChanged(dataTypeFilePath, templateNameDataType, dataType)
	return dataTypeFilename, changed, err
}

func (gen *Generator) outputService(outputPath string, service *Service) (string, bool, error) {
	serviceFilename := service.Name + gen.config.FileExtension
	servicePath := filepath.Join(outputPath, serviceFilename)
	changed, err := gen.outputTemplateIfChanged(servicePath, templateNameService, service)
	return serviceFilename, changed, err
}

func (gen *Generator) outputTemplateIfChanged(filename string, templateName string, data interface{}) (bool, error) {
	// get previous contents of file
	previousOutput := ""
	exists := exists(filename)
	if exists {
		var filebytes []byte
		var err error
		if filebytes, err = ioutil.ReadFile(filename); err != nil {
			return false, err
		}
		previousOutput = string(filebytes)
	}

	// template --> string
	var buffer bytes.Buffer
	if err := gen.templates.ExecuteTemplate(&buffer, templateName, data); err != nil {
		return false, err
	}
	output := buffer.String()

	// no change
	if exists && output == previousOutput {
		return false, nil
	}

	// write to os
	buffer.Reset()
	buffer.WriteString(output)
	if err := ioutil.WriteFile(filename, buffer.Bytes(), os.ModePerm); err != nil {
		return false, err
	}

	return true, nil
}

func exists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

func isValidDirectory(path string) bool {
	if stat, err := os.Stat(path); os.IsNotExist(err) || !stat.IsDir() {
		return false
	}
	return true
}

func joinNamespaceToOutputPath(outputPath string, namespace string) string {
	if filepath.Base(outputPath) == namespace {
		return outputPath
	}

	return filepath.Join(outputPath, namespace)
}
