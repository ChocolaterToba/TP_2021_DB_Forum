package repository

import "dbforum/domain/entity"

type ThreadRepositoryInterface interface {
	CreateThread(thread *entity.Thread) (int, string, error) // Returns threadID, forumname and error if present
	GetThreadByID(threadID int) (*entity.Thread, error)
	GetThreadByThreadname(threadname string) (*entity.Thread, error)
	EditThread(thread *entity.Thread) error
	GetPostsByThreadIDFlat(threadID int, limit int, startAfter int, desc bool) ([]*entity.Post, error)
	VoteThreadByThreadname(threadname string, username string, upvote bool) error
	ChangeVoteThreadByThreadname(threadname string, username string, upvote bool) error
}
