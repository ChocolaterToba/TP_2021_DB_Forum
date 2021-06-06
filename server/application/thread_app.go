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
	CreateThread(thread *entity.Thread) (int, error)
	GetThreadByID(threadID int) (*entity.Thread, error)
	GetThreadByForumname(threadForumname string) (*entity.Thread, error)
	EditThread(thread *entity.Thread) error
	GetPostsByThreadID(threadID int) ([]*entity.Post, error)
	GetPostsByThreadForumname(threadIDstring string) ([]*entity.Post, error)
}

// CreateThread adds new thread to database with passed fields
// It returns thread's assigned ID and nil on success, any number and error on failure
func (threadApp *ThreadApp) CreateThread(thread *entity.Thread) (int, error) {
	return threadApp.threadRepo.CreateThread(thread)
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

// GetThreadByForumname fetches thread with passed thread string ID ("slug") from database
// It returns that thread, nil on success and nil, error on failure
func (threadApp *ThreadApp) GetThreadByForumname(threadForumname string) (*entity.Thread, error) {
	return threadApp.threadRepo.GetThreadByForumname(threadForumname)
}

// GetPostsByThreadID fetches all posts in specified thread
// It returns slice of these posts, nil on success and nil, error on failure
func (threadApp *ThreadApp) GetPostsByThreadID(threadID int) ([]*entity.Post, error) {
	return threadApp.threadRepo.GetPostsByThreadID(threadID)
}

// GetPostsByThreadForumname fetches all posts in specified thread
// It returns slice of these posts, nil on success and nil, error on failure
func (threadApp *ThreadApp) GetPostsByThreadForumname(threadIDstring string) ([]*entity.Post, error) {
	return threadApp.threadRepo.GetPostsByThreadForumname(threadIDstring)
}
