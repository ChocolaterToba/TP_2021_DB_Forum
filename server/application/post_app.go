package application

import (
	"dbforum/domain/entity"
	"dbforum/domain/repository"
	"time"
)

type PostApp struct {
	postRepo  repository.PostRepositoryInterface
	threadApp ThreadAppInterface
}

func NewPostApp(postRepo repository.PostRepositoryInterface, threadApp ThreadAppInterface) *PostApp {
	return &PostApp{
		postRepo:  postRepo,
		threadApp: threadApp,
	}
}

type PostAppInterface interface {
	CreatePost(post *entity.Post) (*entity.Post, error)
	GetPostByID(postID int) (*entity.Post, error)
	GetChildPosts(postID int) ([]*entity.Post, error)
	EditPost(post *entity.Post) error
}

// CreatePost adds new post to database with passed fields
// It returns created port and nil on success, any number and error on failure
func (postApp *PostApp) CreatePost(post *entity.Post) (*entity.Post, error) {
	if post.ParentID != 0 {
		//Checking if parent post exists in same thread
		parentPost, err := postApp.GetPostByID(post.ParentID)
		if err != nil {
			return nil, entity.ParentNotFoundError
		}
		if parentPost.ThreadID != post.ThreadID {
			return nil, entity.ParentNotFoundError
		}
	}

	if post.Created == (time.Time{}) {
		post.Created = time.Now().Truncate(time.Millisecond)
	}

	var err error
	post.PostID, err = postApp.postRepo.CreatePost(post)
	if err != nil {
		return nil, err
	}

	return post, nil
}

// EditPost saves post to database with passed fields
// It returns nil on success and error on failure
func (postApp *PostApp) EditPost(post *entity.Post) error {
	return postApp.postRepo.EditPost(post)
}

// GetPostByID fetches post with passed ID from database
// It returns that post, nil on success and nil, error on failure
func (postApp *PostApp) GetPostByID(postID int) (*entity.Post, error) {
	return postApp.postRepo.GetPostByID(postID)
}

// GetChildPosts fetches posts that are 'children' of passed post
// It returns that post, nil on success and nil, error on failure
func (postApp *PostApp) GetChildPosts(postID int) ([]*entity.Post, error) {
	return postApp.postRepo.GetChildPosts(postID)
}
