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
	CreateForum(forum *entity.ForumCreateInput) (*entity.Forum, error) // Returns created forum, nil on success, conflicting forum, error on failure
	GetForumByID(forumID int) (*entity.Forum, error)
	GetForumByForumname(forumname string) (*entity.Forum, error)
	GetUsersByForumname(forumname string) ([]*entity.User, error)
	GetThreadsByForumname(forumname string) ([]*entity.Thread, error)
}

// CreateForum adds new forum to database with passed fields
// It returns Forum's assigned ID and nil on success, any number and error on failure
func (forumApp *ForumApp) CreateForum(forum *entity.ForumCreateInput) (*entity.Forum, error) {
	createdForum := new(entity.Forum)
	var err error
	createdForum.ForumID, createdForum.Creator, err = forumApp.forumRepo.CreateForum(forum)
	if err != nil {
		switch err {
		case entity.ForumConflictError:
			forumnameConflict, err := forumApp.forumRepo.GetForumByForumname(forum.Forumname)
			switch err {
			case nil:
				return forumnameConflict, entity.ForumConflictError
			case entity.ForumNotFoundError:
				return nil, entity.ForumConflictNotFoundError
			default:
				return nil, err
			}

		default:
			return nil, err
		}
	}

	createdForum.Forumname = forum.Forumname
	createdForum.Title = forum.Title

	return createdForum, nil
}

// GetForumByID fetches forum with passed ID from database
// It returns that Forum, nil on success and nil, error on failure
func (forumApp *ForumApp) GetForumByID(forumID int) (*entity.Forum, error) {
	return forumApp.forumRepo.GetForumByID(forumID)
}

// GetForumByForumname fetches forum with passed id string from database
// It returns that Forum, nil on success and nil, error on failure
func (forumApp *ForumApp) GetForumByForumname(forumname string) (*entity.Forum, error) {
	return forumApp.forumRepo.GetForumByForumname(forumname)
}

// GetUsersByForumname finds all users belonging to specified forum
// It returns slice of them, nil on success and nil, error on failure
func (forumApp *ForumApp) GetUsersByForumname(forumname string) ([]*entity.User, error) {
	return forumApp.forumRepo.GetUsersByForumname(forumname)
}

// GetThreadsByForumname finds all threads belonging to specified forum
// It returns slice of them, nil on success and nil, error on failure
func (forumApp *ForumApp) GetThreadsByForumname(forumname string) ([]*entity.Thread, error) {
	_, err := forumApp.GetForumByForumname(forumname)
	if err != nil {
		return nil, err
	}

	return forumApp.forumRepo.GetThreadsByForumname(forumname)
}
