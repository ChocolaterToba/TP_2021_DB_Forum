package entity

import (
	"strconv"
	"time"

	"github.com/valyala/fasthttp"
)

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

type ThreadGetPostsInput struct {
	Limit      int    `json:"limit"` // JSON is just in case, never actually used
	StartAfter int    `json:"since"`
	SortMode   string `json:"sort"`
	Desc       bool   `json:"desc"`
}

func QueryToThreadGetPostsInput(query *fasthttp.Args) (*ThreadGetPostsInput, error) {
	var err error
	threadInput := new(ThreadGetPostsInput)

	limitString := string(query.Peek("limit"))
	if limitString != "" {
		threadInput.Limit, err = strconv.Atoi(limitString)
		if err != nil {
			return nil, err
		}
	}

	startString := string(query.Peek("since"))
	if startString != "" {
		threadInput.StartAfter, err = strconv.Atoi(startString)
		if err != nil {
			return nil, err
		}
	}

	threadInput.SortMode = string(query.Peek("sort"))
	switch threadInput.SortMode {
	case "flat", "tree", "parent_tree":
		// Intentional no-op
	case "":
		threadInput.SortMode = "flat"
	default:
		return nil, QueryParseError
	}

	descString := string(query.Peek("desc"))
	switch descString {
	case "", "false":
		//Intentionally no-op
	case "true":
		threadInput.Desc = true
	default:
		return nil, QueryParseError
	}

	return threadInput, nil
}
