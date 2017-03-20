package codegen

import (
	"errors"
	"os"
	"path"
	"text/template"
)

const (
	templateNameDataType = "datatype"
	templateNameService  = "service"
)

type Generator struct {
	outputPath string
	config     *GeneratorConfig
	templates  *template.Template
}

type GeneratorConfig struct {
	TemplatesPath string
	FileExtension string
}

func NewGenerator(outputPath string, config *GeneratorConfig) (*Generator, error) {
	// validate output path
	if !isValidDirectory(outputPath) {
		return nil, errors.New("Invalid output path: " + outputPath)
	}

	// validate templates
	if !isValidDirectory(config.TemplatesPath) {
		return nil, errors.New("Invalid templates path: " + config.TemplatesPath)
	}

	templatePattern := path.Join(config.TemplatesPath, "*.gotmpl")
	templates, err := template.ParseGlob(templatePattern)
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

func (gen *Generator) outputNamespace(outputPath string, namespace *Namespace) error {
	namespacePath := path.Join(outputPath, namespace.Name)
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
		err := gen.outputService(namespacePath, service)
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
	dataTypePath := path.Join(outputPath, dataTypeFilename)
	file, err := os.Create(dataTypePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return gen.templates.ExecuteTemplate(file, templateNameDataType, dataType)
}

func (gen *Generator) outputService(outputPath string, service *Service) error {
	serviceFilename := service.Name + gen.config.FileExtension
	servicePath := path.Join(outputPath, serviceFilename)
	file, err := os.Create(servicePath)
	if err != nil {
		return err
	}

	defer file.Close()

	return gen.templates.ExecuteTemplate(file, templateNameService, service)
}
