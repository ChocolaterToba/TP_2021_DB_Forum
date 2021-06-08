package entity

import (
	"strconv"
	"time"

	"github.com/valyala/fasthttp"
)

type Forum struct {
	ForumID      int    `json:"-"`
	Title        string `json:"title"`
	Creator      string `json:"user"`
	Forumname    string `json:"slug"`
	PostsCount   int    `json:"posts"`
	ThreadsCount int    `json:"threads"`
}

type ForumCreateInput struct {
	Title     string `json:"title"`
	Creator   string `json:"user"`
	Forumname string `json:"slug"`
}

type ForumGetUsersInput struct {
	Limit      int    `json:"limit"` // JSON is just in case, never actually used
	StartAfter string `json:"since"`
	Desc       bool   `json:"desc"`
}

func QueryToForumGetUsersInput(query *fasthttp.Args) (*ForumGetUsersInput, error) {
	var err error
	forumInput := new(ForumGetUsersInput)

	limitString := string(query.Peek("limit"))
	if limitString != "" {
		forumInput.Limit, err = strconv.Atoi(limitString)
		if err != nil {
			return nil, err
		}
	}

	forumInput.StartAfter = string(query.Peek("since"))

	descString := string(query.Peek("desc"))
	switch descString {
	case "", "false":
		//Intentionally no-op
	case "true":
		forumInput.Desc = true
	default:
		return nil, QueryParseError
	}

	return forumInput, nil
}

type ForumGetThreadsInput struct {
	Limit     int       `json:"limit"` // JSON is just in case, never actually used
	StartFrom time.Time `json:"since"`
	Desc      bool      `json:"desc"`
}

const inputTimeLayout string = "2006-01-02T15:04:05.000Z"

func QueryToForumGetThreadsInput(query *fasthttp.Args) (*ForumGetThreadsInput, error) {
	var err error
	forumInput := new(ForumGetThreadsInput)

	limitString := string(query.Peek("limit"))
	if limitString != "" {
		forumInput.Limit, err = strconv.Atoi(limitString)
		if err != nil {
			return nil, err
		}
	}

	startString := string(query.Peek("since"))
	if startString != "" {
		forumInput.StartFrom, err = time.Parse(inputTimeLayout, startString)
		if err != nil {
			return nil, err
		}
	}

	descString := string(query.Peek("desc"))
	switch descString {
	case "", "false":
		//Intentionally no-op
	case "true":
		forumInput.Desc = true
	default:
		return nil, QueryParseError
	}

	return forumInput, nil
}
