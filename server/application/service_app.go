package application

import (
	"dbforum/domain/entity"
	"dbforum/domain/repository"
	"sync"
)

type ServiceApp struct {
	serviceRepo repository.ServiceRepositoryInterface
	stats       *entity.StatsInfo
	mu          sync.Mutex
}

func NewServiceApp(serviceRepo repository.ServiceRepositoryInterface) *ServiceApp {
	result := &ServiceApp{
		serviceRepo: serviceRepo,
		stats:       new(entity.StatsInfo),
		mu:          sync.Mutex{},
	}
	result.GetStatsInit()
	return result
}

type ServiceAppInterface interface {
	GetStatsInit() error // Loads statistics from database
	GetStats() (*entity.StatsInfo, error)
	IncrementUsersCount() error
	IncrementForumsCount() error
	IncrementThreadsCount() error
	IncrementPostsCount() error
	TruncateAll() error // Clear entire database
}

// GetStatsInit gets stats from database and copies them to ServiceApp.stats
// It returns nil on success, error on failure
func (serviceApp *ServiceApp) GetStatsInit() error {
	stats, err := serviceApp.serviceRepo.GetStats()
	if err != nil {
		return err
	}

	serviceApp.mu.Lock()
	serviceApp.stats = stats
	serviceApp.mu.Unlock()
	return nil
}

// GetStats returns database statistics: amount of users, forums, thred and posts
// It returns these service, nil on success, nil, error on failure
func (serviceApp *ServiceApp) GetStats() (*entity.StatsInfo, error) {
	stats := new(entity.StatsInfo)
	serviceApp.mu.Lock()
	stats.UsersCount = serviceApp.stats.UsersCount
	stats.ForumsCount = serviceApp.stats.ForumsCount
	stats.ThreadsCount = serviceApp.stats.ThreadsCount
	stats.PostsCount = serviceApp.stats.PostsCount
	serviceApp.mu.Unlock()
	return stats, nil
}

// TruncateAll clears all database tables, removing any user data
func (serviceApp *ServiceApp) TruncateAll() error {
	err := serviceApp.serviceRepo.TruncateAll()
	if err != nil {
		return err
	}

	serviceApp.mu.Lock()
	serviceApp.stats.UsersCount = 0
	serviceApp.stats.ForumsCount = 0
	serviceApp.stats.ThreadsCount = 0
	serviceApp.stats.PostsCount = 0
	serviceApp.mu.Unlock()
	return nil
}

func (serviceApp *ServiceApp) IncrementUsersCount() error {
	serviceApp.mu.Lock()
	serviceApp.stats.UsersCount++
	serviceApp.mu.Unlock()
	return nil
}
func (serviceApp *ServiceApp) IncrementForumsCount() error {
	serviceApp.mu.Lock()
	serviceApp.stats.ForumsCount++
	serviceApp.mu.Unlock()
	return nil
}
func (serviceApp *ServiceApp) IncrementThreadsCount() error {
	serviceApp.mu.Lock()
	serviceApp.stats.ThreadsCount++
	serviceApp.mu.Unlock()
	return nil
}
func (serviceApp *ServiceApp) IncrementPostsCount() error {
	serviceApp.mu.Lock()
	serviceApp.stats.PostsCount++
	serviceApp.mu.Unlock()
	return nil
}
