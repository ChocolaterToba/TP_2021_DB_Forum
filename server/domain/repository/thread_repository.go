package repository

import "dbforum/domain/entity"

type ThreadRepositoryInterface interface {
	CreateThread(thread *entity.Thread) (int, error) // Returns ThreadID and error if present
	GetThreadByID(threadID int) (*entity.Thread, error)
	GetThreadByThreadname(threadname string) (*entity.Thread, error)
	EditThread(thread *entity.Thread) error
	GetPostsByThreadID(threadID int) ([]*entity.Post, error)
	VoteThreadByThreadname(threadname string, username string, upvote bool) error
	ChangeVoteThreadByThreadname(threadname string, username string, upvote bool) error
}
