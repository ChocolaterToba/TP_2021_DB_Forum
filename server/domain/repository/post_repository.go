package repository

import "dbforum/domain/entity"

type PostRepositoryInterface interface {
	CreatePost(post *entity.Post) (int, error)
	GetPostByID(postID int) (*entity.Post, error)
	GetChildPosts(parentID int) ([]*entity.Post, error) // Get immediate children of a post
	EditPost(post *entity.Post) error
}
