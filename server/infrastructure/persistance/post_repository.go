package persistance

import (
	"dbforum/domain/entity"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PostRepo struct {
	postgresDB *pgxpool.Pool
}

func NewPostRepo(postgresDB *pgxpool.Pool) *PostRepo {
	return &PostRepo{postgresDB}
}

func (postRepo *PostRepo) CreatePost(post *entity.Post) (int, error) {
	return -1, nil
}

func (postRepo *PostRepo) GetPostByID(postID int) (*entity.Post, error) {
	return nil, nil
}

func (postRepo *PostRepo) GetChildPosts(postID int) ([]*entity.Post, error) {
	return nil, nil
}

func (postRepo *PostRepo) EditPost(post *entity.Post) error {
	return nil
}
