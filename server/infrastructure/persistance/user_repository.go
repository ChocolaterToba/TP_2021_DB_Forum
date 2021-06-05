package persistance

import (
	"dbforum/domain/entity"

	"github.com/jackc/pgx/v4/pgxpool"
)

type UserRepo struct {
	postgresDB *pgxpool.Pool
}

func NewUserRepo(postgresDB *pgxpool.Pool) *UserRepo {
	return &UserRepo{postgresDB}
}

func (userRepo *UserRepo) CreateUser(user *entity.User) (int, error) {
	return -1, nil
}

func (userRepo *UserRepo) GetUserByID(userID int) (*entity.User, error) {
	return nil, nil
}

func (userRepo *UserRepo) GetUserByUsername(username string) (*entity.User, error) {
	return nil, nil
}

func (userRepo *UserRepo) EditUser(user *entity.User) error {
	return nil
}
