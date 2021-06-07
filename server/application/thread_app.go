package application

import (
	"dbforum/domain/entity"
	"dbforum/domain/repository"
)

type ThreadApp struct {
	threadRepo repository.ThreadRepositoryInterface
}

func NewThreadApp(threadRepo repository.ThreadRepositoryInterface) *ThreadApp {
	return &ThreadApp{threadRepo}
}

type ThreadAppInterface interface {
	CreateThread(thread *entity.Thread) (interface{}, error) // Returns int, nil on success, *entity.Forum, error on failure
	GetThreadByID(threadID int) (*entity.Thread, error)
	GetThreadByThreadname(threadname string) (*entity.Thread, error)
	EditThread(thread *entity.Thread) error
	GetPostsByThreadID(threadID int) ([]*entity.Post, error)
	GetPostsByThreadname(threadname string) ([]*entity.Post, error)
}

// CreateThread adds new thread to database with passed fields
// It returns int, nil on success, *entity.Forum (database conflict), error on failure
func (threadApp *ThreadApp) CreateThread(thread *entity.Thread) (interface{}, error) {
	threadID, err := threadApp.threadRepo.CreateThread(thread)
	if err != nil {
		switch err {
		case entity.ThreadConflictError:
			threadnameConflict, err := threadApp.threadRepo.GetThreadByThreadname(thread.Threadname)
			switch err {
			case nil:
				return threadnameConflict, entity.ForumConflictError
			case entity.ThreadNotFoundError:
				return nil, entity.ThreadConflictNotFoundError
			default:
				return nil, err
			}

		default:
			return nil, err
		}
	}

	return threadID, nil
}

// EditThread saves thread to database with passed fields
// It returns nil on success and error on failure
func (threadApp *ThreadApp) EditThread(thread *entity.Thread) error {
	return threadApp.threadRepo.EditThread(thread)
}

// GetThreadByID fetches thread with passed ID from database
// It returns that thread, nil on success and nil, error on failure
func (threadApp *ThreadApp) GetThreadByID(threadID int) (*entity.Thread, error) {
	return threadApp.threadRepo.GetThreadByID(threadID)
}

// GetThreadByThreadname fetches thread with passed thread string ID ("slug") from database
// It returns that thread, nil on success and nil, error on failure
func (threadApp *ThreadApp) GetThreadByThreadname(threadname string) (*entity.Thread, error) {
	return threadApp.threadRepo.GetThreadByThreadname(threadname)
}

// GetPostsByThreadID fetches all posts in specified thread
// It returns slice of these posts, nil on success and nil, error on failure
func (threadApp *ThreadApp) GetPostsByThreadID(threadID int) ([]*entity.Post, error) {
	return threadApp.threadRepo.GetPostsByThreadID(threadID)
}

// GetPostsByThreadname fetches all posts in specified thread
// It returns slice of these posts, nil on success and nil, error on failure
func (threadApp *ThreadApp) GetPostsByThreadname(threadname string) ([]*entity.Post, error) {
	thread, err := threadApp.GetThreadByThreadname(threadname)
	if err != nil {
		return nil, err
	}

	return threadApp.threadRepo.GetPostsByThreadID(thread.ThreadID)
}
