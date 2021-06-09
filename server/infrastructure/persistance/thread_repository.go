package persistance

import (
	"context"
	"dbforum/domain/entity"
	"strings"

	"github.com/jackc/pgtype"
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
		case strings.Contains(err.Error(), "violates foreign key constraint \"threads_fk_forumname\""):
			return -1, "", entity.ForumNotFoundError
		case strings.Contains(err.Error(), "violates foreign key constraint \"threads_fk_creator\""):
			return -1, "", entity.UserNotFoundError
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

	row := tx.QueryRow(context.Background(), getThreadByIDQuery, threadID)

	thread := entity.Thread{ThreadID: threadID}
	var threadname pgtype.Text
	err = row.Scan(&threadname, &thread.Title, &thread.Creator,
		&thread.Forumname, &thread.Message, &thread.Created, &thread.Rating)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, entity.ThreadNotFoundError
		}
		return nil, err
	}
	if threadname.Status != pgtype.Null {
		thread.Threadname = threadname.String
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

	row := tx.QueryRow(context.Background(), getThreadByThreadnameQuery, threadname)

	thread := entity.Thread{}
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

	posts := make([]*entity.Post, 0)
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

const getPostsByThreadIDTreeQuery string = "SELECT postID, parentID, creator, message, isEdited, created\n" +
	"FROM Posts\n" +
	"WHERE threadID=$1 AND " +
	"COALESCE(" +
	"path>(SELECT path FROM Posts WHERE threadID=$1 AND postID=$2), " +
	"path>=(SELECT path FROM Posts WHERE threadID=$1 AND postID>$2 ORDER BY postID LIMIT 1), " +
	"false)\n" +
	"ORDER BY path\n" +
	"LIMIT $3"

const getPostsByThreadIDTreeDescQuery string = "SELECT postID, parentID, creator, message, isEdited, created\n" +
	"FROM Posts\n" +
	"WHERE threadID=$1 AND " +
	"COALESCE(" +
	"path<(SELECT path FROM Posts WHERE threadID=$1 AND postID=$2), " +
	"path<=(SELECT path FROM Posts WHERE threadID=$1 AND postID<$2 ORDER BY postID DESC LIMIT 1), " +
	"false)\n" +
	"ORDER BY path DESC\n" +
	"LIMIT $3"

const getPostsByThreadIDTreeDescNoStartQuery string = "SELECT postID, parentID, creator, message, isEdited, created\n" +
	"FROM Posts\n" +
	"WHERE threadID=$1\n" +
	"ORDER BY path DESC\n" +
	"LIMIT $2"

func (threadRepo *ThreadRepo) GetPostsByThreadIDTree(threadID int, limit int, startAfter int, desc bool) ([]*entity.Post, error) {
	tx, err := threadRepo.postgresDB.Begin(context.Background())
	if err != nil {
		return nil, entity.TransactionBeginError
	}
	defer tx.Rollback(context.Background())

	var rows pgx.Rows
	switch desc {
	case true:
		switch startAfter {
		case 0:
			rows, err = tx.Query(context.Background(), getPostsByThreadIDTreeDescNoStartQuery, threadID, limit)
		default:
			rows, err = tx.Query(context.Background(), getPostsByThreadIDTreeDescQuery, threadID, startAfter, limit)
		}
	case false:
		rows, err = tx.Query(context.Background(), getPostsByThreadIDTreeQuery, threadID, startAfter, limit)
	}
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, entity.PostNotFoundError
		}
		return nil, err
	}

	posts := make([]*entity.Post, 0)
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

const getPostsByThreadIDTopQuery string = "SELECT postID, parentID, creator, message, isEdited, created\n" +
	"FROM Posts\n" +
	"WHERE threadID=$1 AND " +
	"parentID=0 AND " +
	"COALESCE(" +
	"postID>(SELECT path[1] FROM Posts WHERE threadID=$1 AND postID=$2), " +
	"postID>=(SELECT path[1] FROM Posts WHERE threadID=$1 AND postID>$2 ORDER BY postID LIMIT 1), " +
	"false)\n" +
	"ORDER BY postID\n" +
	"LIMIT $3"

const getPostsByThreadIDTopDescQuery string = "SELECT postID, parentID, creator, message, isEdited, created\n" +
	"FROM Posts\n" +
	"WHERE threadID=$1 AND " +
	"parentID=0 AND " +
	"COALESCE(" +
	"postID<(SELECT path[1] FROM Posts WHERE threadID=$1 AND postID=$2), " +
	"postID<=(SELECT path[1] FROM Posts WHERE threadID=$1 AND postID<$2 ORDER BY postID DESC LIMIT 1), " +
	"false)\n" +
	"ORDER BY postID DESC\n" +
	"LIMIT $3"

const getPostsByThreadIDTopDescNoStartQuery string = "SELECT postID, parentID, creator, message, isEdited, created\n" +
	"FROM Posts\n" +
	"WHERE threadID=$1 AND " +
	"parentID=0\n" +
	"ORDER BY postID DESC\n" +
	"LIMIT $2"

func (threadRepo *ThreadRepo) GetPostsByThreadIDTop(threadID int, limit int, startAfter int, desc bool) ([]*entity.Post, error) {
	tx, err := threadRepo.postgresDB.Begin(context.Background())
	if err != nil {
		return nil, entity.TransactionBeginError
	}
	defer tx.Rollback(context.Background())

	var rows pgx.Rows
	switch desc {
	case true:
		switch startAfter {
		case 0:
			rows, err = tx.Query(context.Background(), getPostsByThreadIDTopDescNoStartQuery, threadID, limit)
		default:
			rows, err = tx.Query(context.Background(), getPostsByThreadIDTopDescQuery, threadID, startAfter, limit)
		}
	case false:
		rows, err = tx.Query(context.Background(), getPostsByThreadIDTopQuery, threadID, startAfter, limit)
	}
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, entity.PostNotFoundError
		}
		return nil, err
	}

	posts := make([]*entity.Post, 0)
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

const insertVoteByThreadnameQuery string = "INSERT INTO Votes (threadID, username, upvote)\n" +
	"VALUES ($1, $2, $3)"
const increaseThreadRatingQuery string = "UPDATE Threads\n" +
	"SET rating = rating + $2\n" +
	"WHERE threadID=$1"
const decreaseThreadRatingQuery string = "UPDATE Threads\n" +
	"SET rating = rating - $2\n" +
	"WHERE threadID=$1"

func (threadRepo *ThreadRepo) VoteThreadByThreadID(threadID int, username string, upvote bool) error {
	tx, err := threadRepo.postgresDB.Begin(context.Background())
	if err != nil {
		return entity.TransactionBeginError
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), insertVoteByThreadnameQuery, threadID, username, upvote)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "violates foreign key constraint \"votes_fk_threadID\""):
			return entity.ThreadNotFoundError
		case strings.Contains(err.Error(), "violates foreign key constraint \"votes_fk_username\""):
			return entity.UserNotFoundError
		case strings.Contains(err.Error(), "constraint"):
			return entity.VoteAlreadyExistsError
		default:
			return err
		}
	}

	switch upvote {
	case true:
		_, err = tx.Exec(context.Background(), increaseThreadRatingQuery, threadID, 1)
		if err != nil {
			return err
		}
	case false:
		_, err = tx.Exec(context.Background(), decreaseThreadRatingQuery, threadID, 1)
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
	"WHERE threadID=$1 AND username=$2"
const updateVoteByThreadnameQuery string = "UPDATE Votes\n" +
	"SET upvote = $3\n" +
	"WHERE threadID=$1 AND username=$2"

func (threadRepo *ThreadRepo) ChangeVoteThreadByThreadID(threadID int, username string, upvote bool) error {
	tx, err := threadRepo.postgresDB.Begin(context.Background())
	if err != nil {
		return entity.TransactionBeginError
	}
	defer tx.Rollback(context.Background())

	row := tx.QueryRow(context.Background(), getVoteByThreadnameQuery, threadID, username)
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
		_, err = tx.Exec(context.Background(), increaseThreadRatingQuery, threadID, 2) // 2 to compensate previous downvote
		if err != nil {
			return err
		}
	case false:
		_, err = tx.Exec(context.Background(), decreaseThreadRatingQuery, threadID, 2)
		if err != nil {
			return err
		}
	}

	_, err = tx.Exec(context.Background(), updateVoteByThreadnameQuery, threadID, username, upvote)
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return entity.TransactionCommitError
	}
	return nil
}
