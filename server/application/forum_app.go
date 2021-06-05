package application

import (
	"dbforum/domain/entity"
	"dbforum/domain/repository"
)

type ForumApp struct {
	forumRepo repository.ForumRepositoryInterface
}

func NewForumApp(forumRepo repository.ForumRepositoryInterface) *ForumApp {
	return &ForumApp{forumRepo}
}

type ForumAppInterface interface {
	CreateForum(forum *entity.Forum) (int, error)
	GetForumByID(forumID int) (*entity.Forum, error)
	GetForumByIDString(forumIDString string) (*entity.Forum, error)
	GetUsersByForumIDString(forumIDString string) ([]*entity.User, error)
	GetThreadsByForumIDString(forumIDString string) ([]*entity.Thread, error)
}

// CreateForum adds new forum to database with passed fields
// It returns Forum's assigned ID and nil on success, any number and error on failure
func (forumApp *ForumApp) CreateForum(forum *entity.Forum) (int, error) {
	return forumApp.forumRepo.CreateForum(forum)
}

// GetForumByID fetches forum with passed ID from database
// It returns that Forum, nil on success and nil, error on failure
func (forumApp *ForumApp) GetForumByID(forumID int) (*entity.Forum, error) {
	return forumApp.forumRepo.GetForumByID(forumID)
}

// GetForumByForumname fetches forum with passed id string from database
// It returns that Forum, nil on success and nil, error on failure
func (forumApp *ForumApp) GetForumByIDString(forumIDString string) (*entity.Forum, error) {
	return forumApp.forumRepo.GetForumByIDString(forumIDString)
}

// GetUsersByForumIDString finds all users belonging to specified forum
// It returns slice of them, nil on success and nil, error on failure
func (forumApp *ForumApp) GetUsersByForumIDString(forumIDString string) ([]*entity.User, error) {
	return forumApp.forumRepo.GetUsersByForumIDString(forumIDString)
}

// GetThreadsByForumIDString finds all threads belonging to specified forum
// It returns slice of them, nil on success and nil, error on failure
func (forumApp *ForumApp) GetThreadsByForumIDString(forumIDString string) ([]*entity.Thread, error) {
	return forumApp.forumRepo.GetThreadsByForumIDString(forumIDString)
}
