package application

import (
	"dbforum/domain/entity"
	"dbforum/domain/repository"
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
	CreatePost(post *entity.Post) (interface{}, error)
	GetPostByID(postID int) (*entity.Post, error)
	GetChildPosts(postID int) ([]*entity.Post, error)
	EditPost(post *entity.Post) error
}

// CreatePost adds new post to database with passed fields
// It returns post's assigned ID and nil on success, any number and error on failure
func (postApp *PostApp) CreatePost(post *entity.Post) (int, error) {
	//Checking if parent post exists in same thread
	posts, err := postApp.threadApp.GetPostsByThreadID(post.ThreadID)
	if err != nil {
		return -1, err
	}
	parentFound := false
	for _, parent := range posts {
		if parent.PostID == post.ParentID {
			parentFound = true
			break
		}
	}
	if !parentFound {
		return -1, entity.ParentNotFoundError
	}

	postID, err := postApp.CreatePost(post)
	if err != nil {
		return -1, err
	}

	return postID, nil
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
