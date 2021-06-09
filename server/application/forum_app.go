package application

import (
	"dbforum/domain/entity"
	"dbforum/domain/repository"
	"time"
)

type ForumApp struct {
	forumRepo  repository.ForumRepositoryInterface
	serviceApp ServiceAppInterface
}

func NewForumApp(forumRepo repository.ForumRepositoryInterface, serviceApp ServiceAppInterface) *ForumApp {
	return &ForumApp{
		forumRepo:  forumRepo,
		serviceApp: serviceApp,
	}
}

type ForumAppInterface interface {
	CreateForum(forum *entity.ForumCreateInput) (*entity.Forum, error) // Returns created forum, nil on success, conflicting forum, error on failure
	GetForumByID(forumID int) (*entity.Forum, error)
	GetForumByForumname(forumname string) (*entity.Forum, error)
	GetUsersByForumname(forumname string, limit int, startAfter string, desc bool) ([]*entity.User, error)
	GetThreadsByForumname(forumname string, limit int, startFrom time.Time, desc bool) ([]*entity.Thread, error)
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

	err = forumApp.serviceApp.IncrementForumsCount()
	if err != nil {
		return nil, err
	}
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

// GetUsersByForumname finds all users belonging to specified forum, ordered, starting after user with username = startAfter
// It returns slice of them, nil on success and nil, error on failure
func (forumApp *ForumApp) GetUsersByForumname(forumname string, limit int, startAfter string, desc bool) ([]*entity.User, error) {
	_, err := forumApp.GetForumByForumname(forumname)
	if err != nil {
		return nil, err
	}

	if limit == 0 {
		limit = 100 // Default limit
	}
	return forumApp.forumRepo.GetUsersByForumname(forumname, limit, startAfter, desc)
}

// GetThreadsByForumname finds all threads belonging to specified forum, ordered, starting with thread with creation date = startFrom
// It returns slice of them, nil on success and nil, error on failure
func (forumApp *ForumApp) GetThreadsByForumname(forumname string, limit int, startFrom time.Time, desc bool) ([]*entity.Thread, error) {
	_, err := forumApp.GetForumByForumname(forumname)
	if err != nil {
		return nil, err
	}

	if limit == 0 {
		limit = 100 // Default limit
	}

	return forumApp.forumRepo.GetThreadsByForumname(forumname, limit, startFrom, desc)
}
