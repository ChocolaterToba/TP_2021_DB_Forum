package persistance

import (
	"context"
	"dbforum/domain/entity"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserRepo struct {
	postgresDB *pgxpool.Pool
}

func NewUserRepo(postgresDB *pgxpool.Pool) *UserRepo {
	return &UserRepo{postgresDB}
}

const createUserQuery string = "INSERT INTO Users (username, email, fullName, description)\n" +
	"values ($1, $2, $3, $4)\n" +
	"RETURNING userID"

func (userRepo *UserRepo) CreateUser(user *entity.User) (int, error) {
	tx, err := userRepo.postgresDB.Begin(context.Background())
	if err != nil {
		return -1, entity.TransactionBeginError
	}
	defer tx.Rollback(context.Background())

	row := tx.QueryRow(context.Background(), createUserQuery,
		user.Username, user.EMail, user.FullName, user.Description)

	newUserID := 0
	err = row.Scan(&newUserID)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			return -1, entity.UserConflictError
		}
		return -1, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return -1, entity.TransactionCommitError
	}

	return newUserID, nil
}

const getUserByIDQuery string = "SELECT username, email, fullName, description\n" +
	"FROM Users WHERE userID=$1"

func (userRepo *UserRepo) GetUserByID(userID int) (*entity.User, error) {
	tx, err := userRepo.postgresDB.Begin(context.Background())
	if err != nil {
		return nil, entity.TransactionBeginError
	}
	defer tx.Rollback(context.Background())

	row := tx.QueryRow(context.Background(), getUserByIDQuery, userID)

	user := entity.User{UserID: userID}
	err = row.Scan(&user.Username, &user.EMail, &user.FullName, &user.Description)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, entity.UserNotFoundError
		}
		return nil, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return nil, entity.TransactionCommitError
	}
	return &user, nil
}

const getUserByUsernameQuery string = "SELECT userID, username, email, fullName, description\n" +
	"FROM Users WHERE username=$1"

func (userRepo *UserRepo) GetUserByUsername(username string) (*entity.User, error) {
	tx, err := userRepo.postgresDB.Begin(context.Background())
	if err != nil {
		return nil, entity.TransactionBeginError
	}
	defer tx.Rollback(context.Background())

	row := tx.QueryRow(context.Background(), getUserByUsernameQuery, username)

	user := entity.User{}
	err = row.Scan(&user.UserID, &user.Username, &user.EMail, &user.FullName, &user.Description)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, entity.UserNotFoundError
		}
		return nil, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return nil, entity.TransactionCommitError
	}
	return &user, nil
}

const getUserByEMailQuery string = "SELECT userID, username, email, fullName, description\n" +
	"FROM Users WHERE email=$1"

func (userRepo *UserRepo) GetUserByEmail(email string) (*entity.User, error) {
	tx, err := userRepo.postgresDB.Begin(context.Background())
	if err != nil {
		return nil, entity.TransactionBeginError
	}
	defer tx.Rollback(context.Background())

	row := tx.QueryRow(context.Background(), getUserByEMailQuery, email)

	user := entity.User{}
	err = row.Scan(&user.UserID, &user.Username, &user.EMail, &user.FullName, &user.Description)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, entity.UserNotFoundError
		}
		return nil, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return nil, entity.TransactionCommitError
	}
	return &user, nil
}

const editUserQuery string = "UPDATE Users\n" +
	"SET email=$2, fullName=$3, description=$4\n" +
	"WHERE username=$1"

func (userRepo *UserRepo) EditUser(user *entity.User) error {
	tx, err := userRepo.postgresDB.Begin(context.Background())
	if err != nil {
		return entity.TransactionBeginError
	}
	defer tx.Rollback(context.Background())

	commandTag, err := tx.Exec(context.Background(), editUserQuery,
		user.Username, user.EMail, user.FullName, user.Description)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			return entity.UserConflictError
		}
		return err
	}

	if commandTag.RowsAffected() != 1 {
		return entity.UserNotFoundError
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return entity.TransactionCommitError
	}

	return nil
}
