package repository

import "dbforum/domain/entity"

type PostRepositoryInterface interface {
	CreatePost(post *entity.Post) (int, error)
	CreatePosts(posts []*entity.Post, forumname string) ([]int, error) // Create posts (must be from same forum)
	GetPostByID(postID int) (*entity.Post, error)
	GetPostTree(parentID int, desc bool) ([]*entity.Post, error) // Get pots's tree - post itself and all of it's children, sorted by tree level
	EditPost(post *entity.Post) error
}
