package repository

import "dbforum/domain/entity"

type ForumRepositoryInterface interface {
	CreateForum(forum *entity.Forum) (int, error)
	GetForumByID(forumID int) (*entity.Forum, error)
	GetForumByForumname(forumname string) (*entity.Forum, error)
	GetUsersByForumname(forumname string) ([]*entity.User, error)
	GetThreadsByForumname(forumname string) ([]*entity.Thread, error)
}
