package repository

import "dbforum/domain/entity"

type ServiceRepositoryInterface interface {
	GetStats() (*entity.StatsInfo, error)
	TruncateAll() error
	VacuumPosts() error
}
