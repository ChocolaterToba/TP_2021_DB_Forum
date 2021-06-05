package repository

import "dbforum/domain/entity"

type ForumRepositoryInterface interface {
	CreateForum(forum *entity.Forum) (int, error)
	GetForumByID(forumID int) (*entity.Forum, error)
	GetForumByIDString(forumIDString string) (*entity.Forum, error)
	GetUsersByForumIDString(forumIDString string) ([]*entity.User, error)
	GetThreadsByForumIDString(forumIDString string) ([]*entity.Thread, error)
}
