package entity

import "time"

type Thread struct {
	ThreadID   int       `json:"id"`
	Threadname string    `json:"slug"`
	Title      string    `json:"title"`
	Creator    string    `json:"author"`
	Forumname  string    `json:"forum"`
	Message    string    `json:"message"`
	Created    time.Time `json:"created"`
	Rating     int       `json:"votes"`
}

type ThreadCreateInput struct {
	Threadname string `json:"slug"`
	Title      string `json:"title"`
	Creator    string `json:"author"`
	Forumname  string `json:"forum"`
	Message    string `json:"message"`
}

type ThreadEditInput struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}
