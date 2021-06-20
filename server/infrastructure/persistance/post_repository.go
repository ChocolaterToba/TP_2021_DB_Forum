package persistance

import (
	"context"
	"dbforum/domain/entity"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PostRepo struct {
	postgresDB *pgxpool.Pool
}

func NewPostRepo(postgresDB *pgxpool.Pool) *PostRepo {
	return &PostRepo{postgresDB}
}

const createPostQuery string = "INSERT INTO Posts (parentID, creator, message, isEdited, threadID, forumname, created, path)\n" +
	"values ($1, $2, $3, $4, $5, $6, $7, " +
	"(SELECT path FROM Posts WHERE postID=$1) || (select currval('posts_postid_seq')::integer))\n" + // Taking parent's path and appending new postID
	"RETURNING postID"
const increaseForumPostCountQuery string = "UPDATE Forums\n" +
	"SET posts_count = posts_count + 1\n" +
	"WHERE forumname=$1"

func (postRepo *PostRepo) CreatePost(post *entity.Post) (int, error) {
	tx, err := postRepo.postgresDB.Begin(context.Background())
	if err != nil {
		return -1, entity.TransactionBeginError
	}
	defer tx.Rollback(context.Background())

	row := tx.QueryRow(context.Background(), createPostQuery,
		post.ParentID, post.Creator, post.Message,
		post.IsEdited, post.ThreadID, post.Forumname, post.Created)

	newPostID := 0
	err = row.Scan(&newPostID)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "violates foreign key constraint \"posts_fk_threadid\""):
			return -1, entity.ThreadNotFoundError
		case strings.Contains(err.Error(), "violates foreign key constraint \"posts_fk_creator\""):
			return -1, entity.UserNotFoundError
		default:
			return -1, err
		}
	}

	_, err = tx.Exec(context.Background(), increaseForumPostCountQuery, post.Forumname)
	if err != nil {
		return -1, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return -1, entity.TransactionCommitError
	}

	return newPostID, nil
}

const getPostByIDQuery string = "SELECT post.parentID, post.creator, post.message, " +
	"thread.forumname, post.isEdited, post.threadID, post.created\n" +
	"FROM Posts as post\n" +
	"INNER JOIN Threads as thread\n" +
	"ON post.threadID = thread.threadID\n" +
	"WHERE post.postID=$1"

func (postRepo *PostRepo) GetPostByID(postID int) (*entity.Post, error) {
	tx, err := postRepo.postgresDB.Begin(context.Background())
	if err != nil {
		return nil, entity.TransactionBeginError
	}
	defer tx.Rollback(context.Background())

	row := tx.QueryRow(context.Background(), getPostByIDQuery, postID)

	post := entity.Post{PostID: postID}
	err = row.Scan(&post.ParentID, &post.Creator, &post.Message, &post.Forumname,
		&post.IsEdited, &post.ThreadID, &post.Created)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, entity.PostNotFoundError
		}
		return nil, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return nil, entity.TransactionCommitError
	}
	return &post, nil
}

const getPostTreeByParentIDQuery string = "SELECT postID, parentID, creator, message, " +
	"isEdited, threadID, created\n" +
	"FROM Posts\n" +
	"WHERE $1=ANY(path)\n" +
	"ORDER BY path"

const getPostTreeByParentIDQueryDesc string = "SELECT postID, parentID, creator, message, " +
	"isEdited, threadID, created\n" +
	"FROM Posts\n" +
	"WHERE $1=ANY(path)\n" +
	"ORDER BY path DESC"

func (postRepo *PostRepo) GetPostTree(parentID int, desc bool) ([]*entity.Post, error) {
	tx, err := postRepo.postgresDB.Begin(context.Background())
	if err != nil {
		return nil, entity.TransactionBeginError
	}
	defer tx.Rollback(context.Background())

	var rows pgx.Rows
	switch desc {
	case true:
		rows, err = tx.Query(context.Background(), getPostTreeByParentIDQueryDesc, parentID)
	case false:
		rows, err = tx.Query(context.Background(), getPostTreeByParentIDQuery, parentID)
	}
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, entity.PostNotFoundError
		}
		return nil, err
	}

	posts := make([]*entity.Post, 0)
	for rows.Next() {
		post := entity.Post{}

		err = rows.Scan(&post.PostID, &post.ParentID, &post.Creator, &post.Message,
			&post.IsEdited, &post.ThreadID, &post.Created)
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

const editPostQuery string = "UPDATE Posts\n" +
	"SET message=$2, isEdited=true\n" +
	"WHERE postID=$1"

func (postRepo *PostRepo) EditPost(post *entity.Post) error {
	tx, err := postRepo.postgresDB.Begin(context.Background())
	if err != nil {
		return entity.TransactionBeginError
	}
	defer tx.Rollback(context.Background())

	commandTag, err := tx.Exec(context.Background(), editPostQuery,
		post.PostID, post.Message)

	if err != nil {
		return err
	}

	if commandTag.RowsAffected() != 1 {
		return entity.PostNotFoundError
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return entity.TransactionCommitError
	}

	return nil
}
