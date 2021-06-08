package repository

import (
	"dbforum/domain/entity"
	"time"
)

type ForumRepositoryInterface interface {
	CreateForum(forum *entity.ForumCreateInput) (int, string, error) // Returns new ForumID, Creator and error if present
	GetForumByID(forumID int) (*entity.Forum, error)
	GetForumByForumname(forumname string) (*entity.Forum, error)
	GetUsersByForumname(forumname string,
		limit int, startAfter string, desc bool) ([]*entity.User, error) // Returns users ordered by username, with username >= (<= if desc) start
	GetThreadsByForumname(forumname string, limit int, startFrom time.Time, desc bool) ([]*entity.Thread, error)
}
