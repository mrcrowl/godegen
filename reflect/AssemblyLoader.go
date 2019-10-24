package reflect

import (
	"educationperfect.com/godegen/cli"
	"path/filepath"
	"strings"
)

type AssemblySet map[string]*Assembly

type AssemblyLoader struct {
	basePath         string
	loadedAssemblies AssemblySet
}

func NewAssemblyLoader(basePath string) *AssemblyLoader {
	return &AssemblyLoader{
		basePath,
		make(AssemblySet),
	}
}

func (loader *AssemblyLoader) Load(assemblyName string) (*Assembly, error) {
	var assembly *Assembly
	var err error
	var alreadyLoaded bool
	var extension = filepath.Ext(assemblyName)
	var baseAssemblyName = assemblyName
	if extension == ".dll" || extension == ".exe" {
		baseAssemblyName = strings.TrimSuffix(assemblyName, extension)
	}

	// already loaded
	if assembly, alreadyLoaded = loader.loadedAssemblies[baseAssemblyName]; !alreadyLoaded {
		assemblyFilepath := filepath.Join(loader.basePath, baseAssemblyName+".dll")
		if assembly, err = loader.loadAssemblyFile(assemblyFilepath); err != nil {
			return nil, err
		}

		loader.loadedAssemblies[baseAssemblyName] = assembly
	}

	return assembly, nil
}

func (loader *AssemblyLoader) loadAssemblyFile(filepath string) (*Assembly, error) {
	var assemblyPEFile *cli.AssemblyPEFile
	var err error

	if assemblyPEFile, err = cli.OpenAssemblyPEFile(filepath); err == nil {
		return &Assembly{assemblyPEFile.Metadata, newTypeCache(), loader}, nil
	}
	return nil, err
}
