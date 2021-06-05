package persistance

import (
	"dbforum/domain/entity"

	"github.com/jackc/pgx/v4/pgxpool"
)

type ForumRepo struct {
	postgresDB *pgxpool.Pool
}

func NewForumRepo(postgresDB *pgxpool.Pool) *ForumRepo {
	return &ForumRepo{postgresDB}
}

func (forumRepo *ForumRepo) CreateForum(forum *entity.Forum) (int, error) {
	return -1, nil
}

func (forumRepo *ForumRepo) GetForumByID(forumID int) (*entity.Forum, error) {
	return nil, nil
}

func (forumRepo *ForumRepo) GetForumByIDString(forumIDString string) (*entity.Forum, error) {
	return nil, nil
}

func (forumRepo *ForumRepo) EditForum(forum *entity.Forum) error {
	return nil
}

func (forumRepo *ForumRepo) GetUsersByForumIDString(forumIDString string) ([]*entity.User, error) {
	return nil, nil
}

func (forumRepo *ForumRepo) GetThreadsByForumIDString(forumIDString string) ([]*entity.Thread, error) {
	return nil, nil
}
