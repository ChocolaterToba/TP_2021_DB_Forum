package persistance

import (
	"context"
	"dbforum/domain/entity"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgtype"
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

//Replacing creator's username to one passed when creating user
const replaceForumCreatorQuery string = "UPDATE Forums\n" +
	"SET creator=username\n" +
	"FROM Users\n" +
	"WHERE forumID=$1 AND creator=username\n" +
	"RETURNING username"

func (forumRepo *ForumRepo) CreateForum(forum *entity.ForumCreateInput) (int, string, error) {
	tx, err := forumRepo.postgresDB.Begin(context.Background())
	if err != nil {
		return -1, "", entity.TransactionBeginError
	}
	defer tx.Rollback(context.Background())

	row := tx.QueryRow(context.Background(), createForumQuery,
		forum.Title, forum.Forumname, forum.Creator)

	newForumID := 0
	err = row.Scan(&newForumID)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate"):
			return -1, "", entity.ForumConflictError
		case strings.Contains(err.Error(), "violates foreign key constraint \"forums_fk_creator\""):
			return -1, "", entity.UserNotFoundError
		default:
			return -1, "", err
		}
	}

	row = tx.QueryRow(context.Background(), replaceForumCreatorQuery, newForumID)
	newCreator := ""
	err = row.Scan(&newCreator)
	if err != nil {
		return -1, "", err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return -1, "", entity.TransactionCommitError
	}

	return newForumID, newCreator, nil
}

const getForumByIDQuery string = "SELECT title, forumname, creator, posts_count, threads_count\n" +
	"FROM Forums WHERE forumID=$1"

func (forumRepo *ForumRepo) GetForumByID(forumID int) (*entity.Forum, error) {
	tx, err := forumRepo.postgresDB.Begin(context.Background())
	if err != nil {
		return nil, entity.TransactionBeginError
	}
	defer tx.Rollback(context.Background())

	row := tx.QueryRow(context.Background(), getForumByIDQuery, forumID)

	forum := entity.Forum{ForumID: forumID}
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

const getForumByForumnameQuery string = "SELECT title, forumname, forumID, creator, posts_count, threads_count\n" +
	"FROM Forums\n" +
	"WHERE forumname=$1"

func (forumRepo *ForumRepo) GetForumByForumname(forumname string) (*entity.Forum, error) {
	tx, err := forumRepo.postgresDB.Begin(context.Background())
	if err != nil {
		return nil, entity.TransactionBeginError
	}
	defer tx.Rollback(context.Background())

	row := tx.QueryRow(context.Background(), getForumByForumnameQuery, forumname)

	forum := entity.Forum{}
	err = row.Scan(&forum.Title, &forum.Forumname, &forum.ForumID, &forum.Creator, &forum.PostsCount, &forum.ThreadsCount)
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

const getUsersByForumnameQuery string = "SELECT DISTINCT ON (users.username)\n" +
	"users.userID, users.username, users.email, users.fullName, users.description\n" +
	"FROM Users as users\n" +
	"LEFT JOIN Posts as post\n" +
	"ON users.username = post.creator and post.forumname = $1\n" +
	"LEFT JOIN Threads as thread\n" +
	"ON users.username = thread.creator and thread.forumname = $1\n" +
	"WHERE (post.forumname IS NOT NULL OR thread.forumname IS NOT NULL) AND " +
	"users.username > $2\n" +
	"ORDER BY users.username\n" +
	"LIMIT $3"

const getUsersByForumnameDescQuery string = "SELECT DISTINCT ON (users.username)\n" +
	"users.userID, users.username, users.email, users.fullName, users.description\n" +
	"FROM Users as users\n" +
	"LEFT JOIN Posts as post\n" +
	"ON users.username = post.creator and post.forumname = $1\n" +
	"LEFT JOIN Threads as thread\n" +
	"ON users.username = thread.creator and thread.forumname = $1\n" +
	"WHERE (post.forumname IS NOT NULL OR thread.forumname IS NOT NULL) AND " +
	"users.username < $2\n" +
	"ORDER BY users.username DESC\n" +
	"LIMIT $3"

const getUsersByForumnameDescNostartQuery string = "SELECT DISTINCT ON (users.username)\n" +
	"users.userID, users.username, users.email, users.fullName, users.description\n" +
	"FROM Users as users\n" +
	"LEFT JOIN Posts as post\n" +
	"ON users.username = post.creator and post.forumname = $1\n" +
	"LEFT JOIN Threads as thread\n" +
	"ON users.username = thread.creator and thread.forumname = $1\n" +
	"WHERE post.forumname IS NOT NULL OR thread.forumname IS NOT NULL\n" +
	"ORDER BY users.username DESC\n" +
	"LIMIT $2"

func (forumRepo *ForumRepo) GetUsersByForumname(forumname string, limit int, startAfter string, desc bool) ([]*entity.User, error) {
	tx, err := forumRepo.postgresDB.Begin(context.Background())
	if err != nil {
		return nil, entity.TransactionBeginError
	}
	defer tx.Rollback(context.Background())

	var rows pgx.Rows
	switch desc {
	case true:
		switch startAfter {
		case "":
			rows, err = tx.Query(context.Background(), getUsersByForumnameDescNostartQuery, forumname, limit)
		default:
			rows, err = tx.Query(context.Background(), getUsersByForumnameDescQuery, forumname, startAfter, limit)
		}
	case false:
		rows, err = tx.Query(context.Background(), getUsersByForumnameQuery, forumname, startAfter, limit)
	}
	if err != nil {
		fmt.Println(err)
		if err == pgx.ErrNoRows {
			return nil, entity.UserNotFoundError
		}
		return nil, err
	}

	users := make([]*entity.User, 0)
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

const getThreadsByForumnameQuery string = "SELECT threadID, threadname, title, creator, forumname, message, created, rating\n" +
	"FROM Threads\n" +
	"WHERE forumname=$1 AND " +
	"created>=$2\n" +
	"ORDER BY created\n" +
	"LIMIT $3"

const getThreadsByForumnameDescQuery string = "SELECT threadID, threadname, title, creator, forumname, message, created, rating\n" +
	"FROM Threads\n" +
	"WHERE forumname=$1 AND " +
	"created<=$2\n" +
	"ORDER BY created DESC\n" +
	"LIMIT $3"

const getThreadsByForumnameDescNostartQuery string = "SELECT threadID, threadname, title, creator, forumname, message, created, rating\n" +
	"FROM Threads\n" +
	"WHERE forumname=$1\n" +
	"ORDER BY created DESC\n" +
	"LIMIT $2"

func (forumRepo *ForumRepo) GetThreadsByForumname(forumname string, limit int, startFrom time.Time, desc bool) ([]*entity.Thread, error) {
	tx, err := forumRepo.postgresDB.Begin(context.Background())
	if err != nil {
		return nil, entity.TransactionBeginError
	}
	defer tx.Rollback(context.Background())

	var rows pgx.Rows
	switch desc {
	case true:
		switch startFrom {
		case time.Time{}:
			rows, err = tx.Query(context.Background(), getThreadsByForumnameDescNostartQuery, forumname, limit)
		default:
			rows, err = tx.Query(context.Background(), getThreadsByForumnameDescQuery, forumname, startFrom, limit)
		}
	case false:
		rows, err = tx.Query(context.Background(), getThreadsByForumnameQuery, forumname, startFrom, limit)
	}
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, entity.ThreadNotFoundError
		}
		return nil, err
	}

	threads := make([]*entity.Thread, 0)
	for rows.Next() {
		thread := entity.Thread{}
		var threadname pgtype.Text
		err = rows.Scan(&thread.ThreadID, &threadname, &thread.Title, &thread.Creator,
			&thread.Forumname, &thread.Message, &thread.Created, &thread.Rating)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		if threadname.Status != pgtype.Null {
			thread.Threadname = threadname.String
		}
		threads = append(threads, &thread)
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return nil, entity.TransactionCommitError
	}
	return threads, nil
}
