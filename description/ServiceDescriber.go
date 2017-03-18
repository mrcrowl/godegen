package description

type ServiceDescriber struct {
}

func NewServiceDescriber() *ServiceDescriber {
	return &ServiceDescriber{}
}

func (res *ServiceDescriber) Describe(serviceTypeName string, assemblyFilepath string) (*ServiceDescription, error) {

	// serviceType := assembly.GetType(serviceTypeName)
	// if serviceType == nil {
	// 	return nil, errors.New("Type '" + serviceTypeName + "' not found in: " + assemblyFilepath)
	// }

	// //TODO
	return nil, nil
}
