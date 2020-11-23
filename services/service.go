package services

type Service interface {
	GetTag() string
	Query(query string) (string, error)
	IsReady() bool
}

type ServiceManager struct {
	Services map[string]Service
}

func NewServiceManager(services ...Service) *ServiceManager {
	serviceManager := ServiceManager{
		Services: make(map[string]Service),
	}
	for _, service := range services {
		serviceManager.Services[service.GetTag()] = service
	}
	return &serviceManager
}
