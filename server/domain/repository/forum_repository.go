package repository

import "dbforum/domain/entity"

type ForumRepositoryInterface interface {
	CreateForum(forum *entity.ForumCreateInput) (int, string, error) // Returns new ForumID, Creator and error if present
	GetForumByID(forumID int) (*entity.Forum, error)
	GetForumByForumname(forumname string) (*entity.Forum, error)
	GetUsersByForumname(forumname string) ([]*entity.User, error)
	GetThreadsByForumname(forumname string) ([]*entity.Thread, error)
}
