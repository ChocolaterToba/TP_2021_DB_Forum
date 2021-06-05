package application

import (
	"dbforum/domain/entity"
	"dbforum/domain/repository"
)

type UserApp struct {
	userRepo repository.UserRepositoryInterface
}

func NewUserApp(userRepo repository.UserRepositoryInterface) *UserApp {
	return &UserApp{userRepo}
}

type UserAppInterface interface {
	CreateUser(user *entity.User) (int, error)
	GetUserByID(userID int) (*entity.User, error)
	GetUserByUsername(username string) (*entity.User, error)
	EditUser(user *entity.User) error
}

// CreateUser adds new user to database with passed fields
// It returns user's assigned ID and nil on success, any number and error on failure
func (userApp *UserApp) CreateUser(user *entity.User) (int, error) {
	return userApp.userRepo.CreateUser(user)
}

// EditUser saves user to database with passed fields
// It returns nil on success and error on failure
func (userApp *UserApp) EditUser(user *entity.User) error {
	return userApp.userRepo.EditUser(user)
}

// GetUserByID fetches user with passed ID from database
// It returns that user, nil on success and nil, error on failure
func (userApp *UserApp) GetUserByID(userID int) (*entity.User, error) {
	return userApp.userRepo.GetUserByID(userID)
}

// GetUserByUsername fetches user with passed username from database
// It returns that user, nil on success and nil, error on failure
func (userApp *UserApp) GetUserByUsername(username string) (*entity.User, error) {
	return userApp.userRepo.GetUserByUsername(username)
}
