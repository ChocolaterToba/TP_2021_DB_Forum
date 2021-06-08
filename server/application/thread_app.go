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
	CreateThread(thread *entity.Thread) (*entity.Thread, error) // Returns created thread, nil on success, conflicting thread, error on failure
	GetThreadByID(threadID int) (*entity.Thread, error)
	GetThreadByThreadname(threadname string) (*entity.Thread, error)
	EditThread(thread *entity.Thread) error
	GetPostsByThreadIDFlat(threadID int, limit int, startAfter int, desc bool) ([]*entity.Post, error)
	GetPostsByThreadnameFlat(threadname string, limit int, startAfter int, desc bool) ([]*entity.Post, error) //[]*entity.Post if mode is flat, [][]*entity.Post if tree/parent_tree
	VoteThreadByID(threadID int, username string, voteAmount int) (*entity.Thread, error)
	VoteThreadByThreadname(threadname string, username string, voteAmount int) (*entity.Thread, error)
}

// CreateThread adds new thread to database with passed fields
// It returns created thread, nil on success, conflicting thread, error on failure
func (threadApp *ThreadApp) CreateThread(thread *entity.Thread) (*entity.Thread, error) {
	threadID, forumname, err := threadApp.threadRepo.CreateThread(thread)
	if err != nil {
		switch err {
		case entity.ThreadConflictError:
			threadnameConflict, err := threadApp.threadRepo.GetThreadByThreadname(thread.Threadname)
			switch err {
			case nil:
				return threadnameConflict, entity.ThreadConflictError
			case entity.ThreadNotFoundError:
				return nil, entity.ThreadConflictNotFoundError
			default:
				return nil, err
			}

		default:
			return nil, err
		}
	}

	thread.ThreadID = threadID
	thread.Forumname = forumname
	return thread, nil
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

// GetPostsByThreadID fetches <limit> posts in specified thread, sorted by creation date, starting after postID = startAfter
// It returns slice of these posts, nil on success and nil, error on failure
func (threadApp *ThreadApp) GetPostsByThreadIDFlat(threadID int, limit int, startAfter int, desc bool) ([]*entity.Post, error) {
	thread, err := threadApp.GetThreadByID(threadID)
	if err != nil {
		return nil, err
	}

	posts, err := threadApp.threadRepo.GetPostsByThreadIDFlat(threadID, limit, startAfter, desc)
	if err != nil {
		return nil, err
	}

	for i := range posts {
		posts[i].Forumname = thread.Forumname
	}
	return posts, nil
}

// GetPostsByThreadname fetches all posts in specified thread
// It returns slice of these posts, nil on success and nil, error on failure
func (threadApp *ThreadApp) GetPostsByThreadnameFlat(threadname string, limit int, startAfter int, desc bool) ([]*entity.Post, error) {
	thread, err := threadApp.GetThreadByThreadname(threadname)
	if err != nil {
		return nil, err
	}

	posts, err := threadApp.threadRepo.GetPostsByThreadIDFlat(thread.ThreadID, limit, startAfter, desc)
	if err != nil {
		return nil, err
	}

	for i := range posts {
		posts[i].Forumname = thread.Forumname
	}
	return posts, nil
}

// VoteThreadByID changes thread's rating, adding vote from username
// It returns voted thread, nil on success and nil, error on failure
func (threadApp *ThreadApp) VoteThreadByID(threadID int, username string, voteAmount int) (*entity.Thread, error) {
	thread, err := threadApp.GetThreadByID(threadID)
	if err != nil {
		return nil, err
	}

	var upvote bool
	switch voteAmount {
	case 1:
		upvote = true
	case -1:
		upvote = false
	default:
		return nil, entity.IncorrectVoteAmountError
	}

	err = threadApp.threadRepo.VoteThreadByThreadname(thread.Threadname, username, upvote)
	if err == entity.VoteAlreadyExistsError {
		err = threadApp.threadRepo.ChangeVoteThreadByThreadname(thread.Threadname, username, upvote)
		if err == entity.VoteAlreadyExistsError {
			return thread, nil
		}

		switch upvote { // Reversing previous vote
		case true:
			thread.Rating++
		case false:
			thread.Rating--
		}
	}

	if err != nil {
		return nil, err
	}

	switch upvote {
	case true:
		thread.Rating++
	case false:
		thread.Rating--
	}

	return thread, nil
}

// VoteThreadByThreadname changes thread's rating, adding vote from username
// It returns voted thread, nil on success and nil, error on failure
func (threadApp *ThreadApp) VoteThreadByThreadname(threadname string, username string, voteAmount int) (*entity.Thread, error) {
	thread, err := threadApp.GetThreadByThreadname(threadname)
	if err != nil {
		return nil, err
	}

	var upvote bool
	switch voteAmount {
	case 1:
		upvote = true
	case -1:
		upvote = false
	default:
		return nil, entity.IncorrectVoteAmountError
	}

	err = threadApp.threadRepo.VoteThreadByThreadname(threadname, username, upvote)
	if err == entity.VoteAlreadyExistsError {
		err = threadApp.threadRepo.ChangeVoteThreadByThreadname(threadname, username, upvote)
		if err == entity.VoteAlreadyExistsError {
			return thread, nil
		}

		switch upvote { // Reversing previous vote
		case true:
			thread.Rating++
		case false:
			thread.Rating--
		}
	}

	if err != nil {
		return nil, err
	}

	switch upvote {
	case true:
		thread.Rating++
	case false:
		thread.Rating--
	}

	return thread, nil
}
