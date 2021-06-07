package entity

import "time"

type Post struct {
	PostID    int       `json:"id"`
	ParentID  int       `json:"parent"`
	Creator   string    `json:"author"`
	Message   string    `json:"message"`
	IsEdited  bool      `json:"isEdited"`
	Forumname string    `json:"forum"`
	ThreadID  int       `json:"thread"`
	Created   time.Time `json:"created"`
}
