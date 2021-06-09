package application

import (
	"dbforum/domain/entity"
	"dbforum/domain/repository"
)

type ServiceApp struct {
	serviceRepo repository.ServiceRepositoryInterface
}

func NewServiceApp(serviceRepo repository.ServiceRepositoryInterface) *ServiceApp {
	return &ServiceApp{serviceRepo}
}

type ServiceAppInterface interface {
	GetStats() (*entity.StatsInfo, error)
	TruncateAll() error // Clear entire database
}

// GetService returns database statistics: amount of users, forums, thred and posts
// It returns these service, nil on success, nil, error on failure
func (serviceApp *ServiceApp) GetStats() (*entity.StatsInfo, error) {
	return serviceApp.serviceRepo.GetStats()
}

// TruncateAll clears all database tables, removing any user data
func (serviceApp *ServiceApp) TruncateAll() error {
	return serviceApp.serviceRepo.TruncateAll()
}
