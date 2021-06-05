package repository

import "dbforum/domain/entity"

type PostRepositoryInterface interface {
	CreatePost(post *entity.Post) (int, error)
	GetPostByID(postID int) (*entity.Post, error)
	GetChildPosts(postID int) ([]*entity.Post, error)
	EditPost(post *entity.Post) error
}
