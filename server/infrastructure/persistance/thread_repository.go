package persistance

import (
	"context"
	"dbforum/domain/entity"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ThreadRepo struct {
	postgresDB *pgxpool.Pool
}

func NewThreadRepo(postgresDB *pgxpool.Pool) *ThreadRepo {
	return &ThreadRepo{postgresDB}
}

const createThreadQuery string = "INSERT INTO Threads (creator, title, forumname, message, created)\n" +
	"values ($1, $2, $3, $4, $5)\n" +
	"RETURNING threadID"
const createThreadWithThreadnameQuery string = "INSERT INTO Threads (threadname, creator, title, forumname, message, created)\n" +
	"values ($1, $2, $3, $4, $5, $6)\n" +
	"RETURNING threadID"

//Replacing thread's forumname to one passed when creating forumname
const replaceThreadForumnameQuery string = "UPDATE Threads as thread\n" +
	"SET forumname=forum.forumname\n" +
	"FROM Forums as forum\n" +
	"WHERE thread.threadID=$1 AND thread.forumname=forum.forumname\n" +
	"RETURNING forum.forumname"
const increaseForumThreadCountQuery string = "UPDATE Forums\n" +
	"SET threads_count = threads_count + 1\n" +
	"WHERE forumname=$1"

func (threadRepo *ThreadRepo) CreateThread(thread *entity.Thread) (int, string, error) {
	tx, err := threadRepo.postgresDB.Begin(context.Background())
	if err != nil {
		return -1, "", entity.TransactionBeginError
	}
	defer tx.Rollback(context.Background())

	var row pgx.Row
	switch thread.Threadname {
	case "":
		row = tx.QueryRow(context.Background(), createThreadQuery,
			thread.Creator, thread.Title, thread.Forumname, thread.Message, thread.Created)
	default:
		row = tx.QueryRow(context.Background(), createThreadWithThreadnameQuery,
			thread.Threadname, thread.Creator, thread.Title, thread.Forumname, thread.Message, thread.Created)
	}

	newThreadID := 0
	err = row.Scan(&newThreadID)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate"):
			return -1, "", entity.ThreadConflictError
		case strings.Contains(err.Error(), "violates foreign key"):
			return -1, "", entity.ForumNotFoundError // TODO: differentiate between user not found and forum not found
		default:
			return -1, "", err
		}
	}

	row = tx.QueryRow(context.Background(), replaceThreadForumnameQuery, newThreadID)
	newForumname := ""
	err = row.Scan(&newForumname)
	if err != nil {
		return -1, "", err
	}

	_, err = tx.Exec(context.Background(), increaseForumThreadCountQuery, thread.Forumname)
	if err != nil {
		return -1, "", err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return -1, "", entity.TransactionCommitError
	}

	return newThreadID, newForumname, nil
}

const getThreadByIDQuery string = "SELECT threadname, title, creator, forumname, message, created, rating\n" +
	"FROM Threads WHERE threadID=$1"

func (threadRepo *ThreadRepo) GetThreadByID(threadID int) (*entity.Thread, error) {
	tx, err := threadRepo.postgresDB.Begin(context.Background())
	if err != nil {
		return nil, entity.TransactionBeginError
	}
	defer tx.Rollback(context.Background())

	thread := entity.Thread{ThreadID: threadID}

	row := tx.QueryRow(context.Background(), getThreadByIDQuery, threadID)
	err = row.Scan(&thread.Threadname, &thread.Title, &thread.Creator,
		&thread.Forumname, &thread.Message, &thread.Created, &thread.Rating)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, entity.ThreadNotFoundError
		}
		return nil, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return nil, entity.TransactionCommitError
	}
	return &thread, nil
}

const getThreadByThreadnameQuery string = "SELECT threadID, threadname, title, creator, forumname, message, created, rating\n" +
	"FROM Threads WHERE threadname=$1"

func (threadRepo *ThreadRepo) GetThreadByThreadname(threadname string) (*entity.Thread, error) {
	tx, err := threadRepo.postgresDB.Begin(context.Background())
	if err != nil {
		return nil, entity.TransactionBeginError
	}
	defer tx.Rollback(context.Background())

	thread := entity.Thread{}

	row := tx.QueryRow(context.Background(), getThreadByThreadnameQuery, threadname)
	err = row.Scan(&thread.ThreadID, &thread.Threadname, &thread.Title, &thread.Creator,
		&thread.Forumname, &thread.Message, &thread.Created, &thread.Rating)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, entity.ThreadNotFoundError
		}
		return nil, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return nil, entity.TransactionCommitError
	}
	return &thread, nil
}

const editThreadQuery string = "UPDATE Threads\n" +
	"SET title=$2, message=$3\n" +
	"WHERE threadID=$1"

func (threadRepo *ThreadRepo) EditThread(thread *entity.Thread) error {
	tx, err := threadRepo.postgresDB.Begin(context.Background())
	if err != nil {
		return entity.TransactionBeginError
	}
	defer tx.Rollback(context.Background())

	commandTag, err := tx.Exec(context.Background(), editThreadQuery,
		thread.ThreadID, thread.Title, thread.Message)

	if err != nil {
		return err
	}

	if commandTag.RowsAffected() != 1 {
		return entity.ThreadNotFoundError
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return entity.TransactionCommitError
	}

	return nil
}

const getPostsByThreadIDFlatQuery string = "SELECT postID, parentID, creator, message, isEdited, created\n" +
	"FROM Posts\n" +
	"WHERE threadID=$1 AND " +
	"postID>$2\n" +
	"ORDER BY postID, created\n" +
	"LIMIT $3"

const getPostsByThreadIDFlatDescQuery string = "SELECT postID, parentID, creator, message, isEdited, created\n" +
	"FROM Posts\n" +
	"WHERE threadID=$1 AND " +
	"postID<$2\n" +
	"ORDER BY postID DESC, created DESC\n" +
	"LIMIT $3"

const getPostsByThreadIDFlatDescNoStartQuery string = "SELECT postID, parentID, creator, message, isEdited, created\n" +
	"FROM Posts\n" +
	"WHERE threadID=$1\n" +
	"ORDER BY postID DESC, created DESC\n" +
	"LIMIT $2"

func (threadRepo *ThreadRepo) GetPostsByThreadIDFlat(threadID int, limit int, startAfter int, desc bool) ([]*entity.Post, error) {
	tx, err := threadRepo.postgresDB.Begin(context.Background())
	if err != nil {
		return nil, entity.TransactionBeginError
	}
	defer tx.Rollback(context.Background())

	posts := make([]*entity.Post, 0)
	var rows pgx.Rows
	switch desc {
	case true:
		switch startAfter {
		case 0:
			rows, err = tx.Query(context.Background(), getPostsByThreadIDFlatDescNoStartQuery, threadID, limit)
		default:
			rows, err = tx.Query(context.Background(), getPostsByThreadIDFlatDescQuery, threadID, startAfter, limit)
		}
	case false:
		rows, err = tx.Query(context.Background(), getPostsByThreadIDFlatQuery, threadID, startAfter, limit)
	}

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, entity.PostNotFoundError
		}
		return nil, err
	}

	for rows.Next() {
		post := entity.Post{ThreadID: threadID}

		err = rows.Scan(&post.PostID, &post.ParentID, &post.Creator, &post.Message, &post.IsEdited, &post.Created)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return nil, entity.TransactionCommitError
	}
	return posts, nil
}

const insertVoteByThreadnameQuery string = "INSERT INTO Votes (threadname, username, upvote)\n" +
	"VALUES ($1, $2, $3)"
const increaseThreadRatingQuery string = "UPDATE Threads\n" +
	"SET rating = rating + $2\n" +
	"WHERE threadname=$1"
const decreaseThreadRatingQuery string = "UPDATE Threads\n" +
	"SET rating = rating - $2\n" +
	"WHERE threadname=$1"

func (threadRepo *ThreadRepo) VoteThreadByThreadname(threadname string, username string, upvote bool) error {
	tx, err := threadRepo.postgresDB.Begin(context.Background())
	if err != nil {
		return entity.TransactionBeginError
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), insertVoteByThreadnameQuery, threadname, username, upvote)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "violates foreign key"):
			return entity.ThreadNotFoundError // TODO: differentiate between user not found and thread not found
		case strings.Contains(err.Error(), "constraint"):
			return entity.VoteAlreadyExistsError
		default:
			return err
		}
	}

	switch upvote {
	case true:
		_, err = tx.Exec(context.Background(), increaseThreadRatingQuery, threadname, 1)
		if err != nil {
			return err
		}
	case false:
		_, err = tx.Exec(context.Background(), decreaseThreadRatingQuery, threadname, 1)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return entity.TransactionCommitError
	}
	return nil
}

const getVoteByThreadnameQuery string = "SELECT upvote\n" +
	"FROM Votes\n" +
	"WHERE threadname=$1 AND username=$2"
const updateVoteByThreadnameQuery string = "UPDATE Votes\n" +
	"SET upvote = $3\n" +
	"WHERE threadname=$1 AND username=$2"

func (threadRepo *ThreadRepo) ChangeVoteThreadByThreadname(threadname string, username string, upvote bool) error {
	tx, err := threadRepo.postgresDB.Begin(context.Background())
	if err != nil {
		return entity.TransactionBeginError
	}
	defer tx.Rollback(context.Background())

	row := tx.QueryRow(context.Background(), getVoteByThreadnameQuery, threadname, username)
	var wasUpvoted bool
	err = row.Scan(&wasUpvoted)
	if err != nil {
		if err == pgx.ErrNoRows {
			return entity.VoteNotFoundError
		}

		return err
	}

	if wasUpvoted == upvote {
		return entity.VoteAlreadyExistsError
	}

	switch upvote {
	case true:
		_, err = tx.Exec(context.Background(), increaseThreadRatingQuery, threadname, 2) // 2 to compensate previous downvote
		if err != nil {
			return err
		}
	case false:
		_, err = tx.Exec(context.Background(), decreaseThreadRatingQuery, threadname, 2)
		if err != nil {
			return err
		}
	}

	_, err = tx.Exec(context.Background(), updateVoteByThreadnameQuery, threadname, username, upvote)
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return entity.TransactionCommitError
	}
	return nil
}
