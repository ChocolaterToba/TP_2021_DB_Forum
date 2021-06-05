package repository

import "dbforum/domain/entity"

type UserRepositoryInterface interface {
	CreateUser(user *entity.User) (int, error)
	GetUserByID(userID int) (*entity.User, error)
	GetUserByUsername(username string) (*entity.User, error)
	EditUser(user *entity.User) error
}
