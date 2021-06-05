package entity

import "time"

type Thread struct {
	ThreadID       int       `json:"id"`
	Title          string    `json:"title"`
	AuthorUsername string    `json:"author"`
	ForumName      string    `json:"forum"`
	Message        string    `json:"message"`
	ForumIDString  string    `json:"slug"`
	Created        time.Time `json:"created"`
	Rating         int       `json:"votes"`
}
