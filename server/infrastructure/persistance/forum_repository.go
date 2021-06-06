package persistance

import (
	"context"
	"dbforum/domain/entity"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ForumRepo struct {
	postgresDB *pgxpool.Pool
}

func NewForumRepo(postgresDB *pgxpool.Pool) *ForumRepo {
	return &ForumRepo{postgresDB}
}

const createForumQuery string = "INSERT INTO Forums (title, forumname, creator)\n" +
	"values ($1, $2, $3)\n" +
	"RETURNING forumID"

func (forumRepo *ForumRepo) CreateForum(forum *entity.Forum) (int, error) {
	tx, err := forumRepo.postgresDB.Begin(context.Background())
	if err != nil {
		return -1, entity.TransactionBeginError
	}
	defer tx.Rollback(context.Background())

	row := tx.QueryRow(context.Background(), createForumQuery,
		forum.Title, forum.Forumname, forum.Creator)

	newForumID := 0
	err = row.Scan(&newForumID)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate"):
			return -1, entity.ForumConflictError
		case strings.Contains(err.Error(), "violates foreighn key"):
			return -1, entity.UserNotFoundError
		default:
			return -1, err
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return -1, entity.TransactionCommitError
	}

	return newForumID, nil
}

const getForumByIDQuery string = "SELECT title, forumname, creator, posts_count, threads_count\n" +
	"FROM Forums WHERE forumID=$1"

func (forumRepo *ForumRepo) GetForumByID(forumID int) (*entity.Forum, error) {
	tx, err := forumRepo.postgresDB.Begin(context.Background())
	if err != nil {
		return nil, entity.TransactionBeginError
	}
	defer tx.Rollback(context.Background())

	forum := entity.Forum{ForumID: forumID}

	row := tx.QueryRow(context.Background(), getUserByIDQuery, forumID)
	err = row.Scan(&forum.Title, &forum.Forumname, &forum.Creator, &forum.PostsCount, &forum.ThreadsCount)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, entity.ForumNotFoundError
		}
		return nil, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return nil, entity.TransactionCommitError
	}
	return &forum, nil
}

const getForumByForumnameQuery string = "SELECT title, forumID, creator, posts_count, threads_count\n" +
	"FROM Forums WHERE forumname=$1"

func (forumRepo *ForumRepo) GetForumByForumname(forumname string) (*entity.Forum, error) {
	tx, err := forumRepo.postgresDB.Begin(context.Background())
	if err != nil {
		return nil, entity.TransactionBeginError
	}
	defer tx.Rollback(context.Background())

	forum := entity.Forum{Forumname: forumname}

	row := tx.QueryRow(context.Background(), getUserByIDQuery, forumname)
	err = row.Scan(&forum.Title, &forum.ForumID, &forum.Creator, &forum.PostsCount, &forum.ThreadsCount)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, entity.ForumNotFoundError
		}
		return nil, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return nil, entity.TransactionCommitError
	}
	return &forum, nil
}

const getUsersByForumnameQuery string = "SELECT user.userID, user.username, user.email, user.fullName, user.description\n" +
	"FROM Users as user\n" +
	"INNER JOIN POSTS as post\n" +
	"ON post.authorUsername = user.username\n" +
	"WHERE post.forumname=$1"

func (forumRepo *ForumRepo) GetUsersByForumname(forumname string) ([]*entity.User, error) {
	tx, err := forumRepo.postgresDB.Begin(context.Background())
	if err != nil {
		return nil, entity.TransactionBeginError
	}
	defer tx.Rollback(context.Background())

	users := make([]*entity.User, 0)
	rows, err := tx.Query(context.Background(), getUsersByForumnameQuery, forumname)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, entity.UserNotFoundError
		}
		return nil, err
	}

	for rows.Next() {
		user := entity.User{}

		err = rows.Scan(&user.UserID, &user.Username, &user.EMail, &user.FullName, &user.Description)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return nil, entity.TransactionCommitError
	}
	return users, nil
}

const getThreadsByForumnameQuery string = "SELECT threadID, title, authorUsername, forumName, message, created, rating\n" +
	"FROM Threads\n" +
	"WHERE forumname=$1"

func (forumRepo *ForumRepo) GetThreadsByForumname(forumname string) ([]*entity.Thread, error) {
	tx, err := forumRepo.postgresDB.Begin(context.Background())
	if err != nil {
		return nil, entity.TransactionBeginError
	}
	defer tx.Rollback(context.Background())

	threads := make([]*entity.Thread, 0)
	rows, err := tx.Query(context.Background(), getThreadsByForumnameQuery, forumname)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, entity.ThreadNotFoundError
		}
		return nil, err
	}

	for rows.Next() {
		thread := entity.Thread{}

		err = rows.Scan(&thread.ThreadID, &thread.Title, &thread.AuthorUsername,
			&thread.Forumname, &thread.Message, &thread.Created, &thread.Rating)
		if err != nil {
			return nil, err
		}
		threads = append(threads, &thread)
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return nil, entity.TransactionCommitError
	}
	return threads, nil
}
