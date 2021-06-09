package application

import (
	"dbforum/domain/entity"
	"dbforum/domain/repository"
)

type UserApp struct {
	userRepo   repository.UserRepositoryInterface
	serviceApp ServiceAppInterface
}

func NewUserApp(userRepo repository.UserRepositoryInterface, serviceApp ServiceAppInterface) *UserApp {
	return &UserApp{
		userRepo:   userRepo,
		serviceApp: serviceApp,
	}
}

type UserAppInterface interface {
	CreateUser(user *entity.User) (interface{}, error) // Returns int, nil on success, []*entity.User, error on failure
	GetUserByID(userID int) (*entity.User, error)
	GetUserByUsername(username string) (*entity.User, error)
	EditUser(user *entity.User) error
}

// CreateUser adds new user to database with passed fields
// It returns int, nil on success, {[]*entity.User (slice of conflicting users) OR nil}, error on failure
func (userApp *UserApp) CreateUser(user *entity.User) (interface{}, error) {
	userID, err := userApp.userRepo.CreateUser(user)
	if err != nil {
		switch err {
		case entity.UserConflictError:
			var users []*entity.User
			usernameConflict, err := userApp.userRepo.GetUserByUsername(user.Username)
			switch err {
			case nil:
				users = append(users, usernameConflict)
			case entity.UserNotFoundError:
				// Intentionally a no-op
			default:
				return nil, err
			}

			emailConflict, err := userApp.userRepo.GetUserByEmail(user.EMail)
			switch err {
			case nil:
				users = append(users, emailConflict)
			case entity.UserNotFoundError:
				// Intentionally a no-op
			default:
				return nil, err
			}

			if len(users) == 0 {
				return nil, entity.UserConflictNotFoundError
			}

			if len(users) == 2 && users[0].UserID == users[1].UserID { // Both conflicts are the same user
				users = users[:1]
			}

			return users, entity.UserConflictError

		default:
			return nil, err
		}
	}

	err = userApp.serviceApp.IncrementUsersCount()
	if err != nil {
		return nil, err
	}
	return userID, nil
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
