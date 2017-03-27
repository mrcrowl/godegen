package codegen

import (
	"errors"
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

type Generator struct {
	outputPath string
	config     *GeneratorConfig
	templates  *template.Template
}

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

func isValidDirectory(path string) bool {
	if stat, err := os.Stat(path); os.IsNotExist(err) || !stat.IsDir() {
		return false
	}
	return true
}

func (gen *Generator) OutputServiceDescription(descr *ServiceDescription) error {
	for _, namespace := range descr.Namespaces {
		err := gen.outputNamespace(gen.outputPath, namespace)
		if err != nil {
			return err
		}
	}

	return nil
}

func joinNamespaceToOutputPath(outputPath string, namespace string) string {
	if filepath.Base(outputPath) == namespace {
		return outputPath
	}

	return filepath.Join(outputPath, namespace)
}

func (gen *Generator) outputNamespace(outputPath string, namespace *Namespace) error {
	namespacePath := joinNamespaceToOutputPath(outputPath, namespace.Name)
	os.Mkdir(namespacePath, os.ModePerm)

	// datatype
	for _, dataType := range namespace.DataTypes {
		err := gen.outputDataType(namespacePath, dataType)
		if err != nil {
			return err
		}
	}

	// service
	for _, service := range namespace.Services {
		var servicePath string
		if gen.config.KeepServicesInNamespace {
			servicePath = namespacePath
		} else {
			servicePath = gen.config.OutputPath
		}
		err := gen.outputService(servicePath, service)
		if err != nil {
			return err
		}
	}

	// sub-namespaces
	for _, namespace := range namespace.Namespaces {
		err := gen.outputNamespace(namespacePath, namespace)
		if err != nil {
			return err
		}
	}

	return nil
}

func (gen *Generator) outputDataType(outputPath string, dataType *DataType) error {
	dataTypeFilename := dataType.Name + gen.config.FileExtension
	dataTypePath := filepath.Join(outputPath, dataTypeFilename)
	file, err := os.Create(dataTypePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return gen.templates.ExecuteTemplate(file, templateNameDataType, dataType)
}

func (gen *Generator) outputService(outputPath string, service *Service) error {
	serviceFilename := service.Name + gen.config.FileExtension
	servicePath := filepath.Join(outputPath, serviceFilename)
	file, err := os.Create(servicePath)
	if err != nil {
		return err
	}

	defer file.Close()

	return gen.templates.ExecuteTemplate(file, templateNameService, service)
}
