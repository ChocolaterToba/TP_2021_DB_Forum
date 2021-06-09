package repository

import "dbforum/domain/entity"

type ThreadRepositoryInterface interface {
	CreateThread(thread *entity.Thread) (int, string, error) // Returns threadID, forumname and error if present
	GetThreadByID(threadID int) (*entity.Thread, error)
	GetThreadByThreadname(threadname string) (*entity.Thread, error)
	EditThread(thread *entity.Thread) error
	GetPostsByThreadIDFlat(threadID int, limit int, startAfter int, desc bool) ([]*entity.Post, error)
	GetPostsByThreadIDTree(threadID int, limit int, startAfter int, desc bool) ([]*entity.Post, error)
	GetPostsByThreadIDParentTree(threadID int, limit int, startAfter int, desc bool) ([]*entity.Post, error) // Returns only top posts with no parent
	VoteThreadByThreadID(threadID int, username string, upvote bool) error
	ChangeVoteThreadByThreadID(threadID int, username string, upvote bool) error
}
