package repository

import "dbforum/domain/entity"

type ThreadRepositoryInterface interface {
	CreateThread(thread *entity.Thread) (int, error)
	GetThreadByID(threadID int) (*entity.Thread, error)
	GetThreadByForumname(threadForumname string) (*entity.Thread, error)
	EditThread(thread *entity.Thread) error
	GetPostsByThreadID(threadID int) ([]*entity.Post, error)
	GetPostsByThreadForumname(threadIDstring string) ([]*entity.Post, error)
}
