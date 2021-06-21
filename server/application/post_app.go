package application

import (
	"dbforum/domain/entity"
	"dbforum/domain/repository"
	"time"
)

type PostApp struct {
	postRepo   repository.PostRepositoryInterface
	userRepo   repository.UserRepositoryInterface
	threadRepo repository.ThreadRepositoryInterface
	forumRepo  repository.ForumRepositoryInterface
	serviceApp ServiceAppInterface
}

func NewPostApp(postRepo repository.PostRepositoryInterface, userRepo repository.UserRepositoryInterface,
	threadRepo repository.ThreadRepositoryInterface, forumRepo repository.ForumRepositoryInterface,
	serviceApp ServiceAppInterface) *PostApp {
	return &PostApp{
		postRepo:   postRepo,
		userRepo:   userRepo,
		threadRepo: threadRepo,
		forumRepo:  forumRepo,
		serviceApp: serviceApp,
	}
}

type PostAppInterface interface {
	CreatePost(post *entity.Post) (*entity.Post, error)
	CreatePosts(posts []*entity.Post, thread *entity.Thread) ([]*entity.Post, error)
	GetPostByID(postID int) (*entity.Post, error)
	GetPostRelatives(post *entity.Post, relatives map[string]bool) (*entity.PostFullOutput, error)
	GetPostTree(postID int, desc bool) ([]*entity.Post, error)
	EditPost(post *entity.Post) error
}

// CreatePost adds new post to database with passed fields
// It returns created post and nil on success, nil and error on failure
func (postApp *PostApp) CreatePost(post *entity.Post) (*entity.Post, error) {
	if post.ParentID != 0 {
		//Checking if parent post exists in same thread
		parentPost, err := postApp.GetPostByID(post.ParentID)
		if err != nil {
			return nil, entity.ParentNotFoundError
		}
		if parentPost.ThreadID != post.ThreadID {
			return nil, entity.ParentInAnotherThreadError
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

	err = postApp.serviceApp.IncrementPostsCount(1)
	if err != nil {
		return nil, err
	}
	return post, nil
}

// CreatePosts adds new posts from passed thread to database
// It returns created posts and nil on success, nil and error on failure
func (postApp *PostApp) CreatePosts(posts []*entity.Post, thread *entity.Thread) ([]*entity.Post, error) {
	if len(posts) == 0 {
		return posts, nil
	}

	var err error

	creationDate := time.Now().Truncate(time.Millisecond) // Creation time is same for all posts where it's not specified
	for i, post := range posts {
		posts[i].ThreadID = thread.ThreadID
		posts[i].Forumname = thread.Forumname
		if posts[i].Created == (time.Time{}) {
			posts[i].Created = creationDate
		}

		if post.ParentID != 0 {
			//Checking if parent post exists in same thread
			parentPost, err := postApp.GetPostByID(post.ParentID)
			if err != nil {
				return nil, entity.ParentNotFoundError
			}
			if parentPost.ThreadID != post.ThreadID {
				return nil, entity.ParentInAnotherThreadError
			}
		}
	}

	postIDs, err := postApp.postRepo.CreatePosts(posts)
	if err != nil || len(postIDs) != len(posts) {
		return nil, err
	}

	for i, id := range postIDs {
		posts[i].PostID = id
	}

	err = postApp.serviceApp.IncrementPostsCount(len(posts))
	if err != nil {
		return nil, err
	}
	return posts, nil
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

// GetPostRelatives fetches post's thread, forum and user if specified
// It returns struct with fetched post, thread..., nil on success and nil, error on failure
func (postApp *PostApp) GetPostRelatives(post *entity.Post, relatives map[string]bool) (*entity.PostFullOutput, error) {
	output := new(entity.PostFullOutput)
	output.PostOutput = post

	var err error
	if relatives["user"] {
		output.UserOutput, err = postApp.userRepo.GetUserByUsername(post.Creator)
		if err != nil {
			return nil, err
		}
	}
	if relatives["thread"] || relatives["forum"] {
		thread, err := postApp.threadRepo.GetThreadByID(post.ThreadID)
		if err != nil {
			return nil, err
		}

		if relatives["thread"] {
			output.ThreadOutput = thread
		}

		if relatives["forum"] {
			output.ForumOutput, err = postApp.forumRepo.GetForumByForumname(thread.Forumname)
			if err != nil {
				return nil, err
			}
		}
	}

	return output, nil
}

// GetPostTree fetches post, it's children, their children... in tree order
// It returns that post, nil on success and nil, error on failure
func (postApp *PostApp) GetPostTree(postID int, desc bool) ([]*entity.Post, error) {
	return postApp.postRepo.GetPostTree(postID, desc)
}
