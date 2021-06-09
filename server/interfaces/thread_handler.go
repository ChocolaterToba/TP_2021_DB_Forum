package interfaces

import (
	"dbforum/application"
	"dbforum/domain/entity"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/valyala/fasthttp"
)

type ThreadInfo struct {
	threadApp application.ThreadAppInterface
}

func NewThreadInfo(threadApp application.ThreadAppInterface) *ThreadInfo {
	return &ThreadInfo{threadApp: threadApp}
}

func (threadInfo *ThreadInfo) CreateThread(ctx *fasthttp.RequestCtx) {
	forumnameInterface := ctx.UserValue("forumname")
	threadInput := new(entity.Thread)

	switch forumnameInterface.(type) {
	case string:
		threadInput.Forumname = forumnameInterface.(string)
	default:
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	err := json.Unmarshal(ctx.Request.Body(), threadInput)
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	//TODO: validate

	if threadInput.Created == (time.Time{}) {
		threadInput.Created = time.Now().Truncate(time.Millisecond)
	}

	newThread, err := threadInfo.threadApp.CreateThread(threadInput)
	if err != nil {
		switch err {
		case entity.ThreadConflictError:
			responseBody, err := json.Marshal(newThread)
			if err != nil {
				ctx.SetStatusCode(http.StatusInternalServerError)
				return
			}

			ctx.SetStatusCode(http.StatusConflict)
			ctx.SetContentType("application/json")
			ctx.SetBody(responseBody)
			return
			//TODO: rework
		case entity.ForumNotFoundError:
			responseBody, err := json.Marshal(entity.MessageOutput{"Can't find forum"})
			if err != nil {
				ctx.SetStatusCode(http.StatusInternalServerError)
				return
			}

			ctx.SetStatusCode(http.StatusNotFound)
			ctx.SetContentType("application/json")
			ctx.SetBody(responseBody)
			return
		case entity.UserNotFoundError:
			responseBody, err := json.Marshal(entity.MessageOutput{"Could not find user"})
			if err != nil {
				ctx.SetStatusCode(http.StatusInternalServerError)
				return
			}

			ctx.SetStatusCode(http.StatusNotFound)
			ctx.SetContentType("application/json")
			ctx.SetBody(responseBody)
			return
		default:
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}
	}

	responseBody, err := json.Marshal(newThread)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(http.StatusCreated)
	ctx.SetContentType("application/json")
	ctx.SetBody(responseBody)
}

func (threadInfo *ThreadInfo) GetThread(ctx *fasthttp.RequestCtx) {
	threadnameInterface := ctx.UserValue("threadnameOrID")
	var threadname string

	switch threadnameInterface.(type) {
	case string:
		threadname = threadnameInterface.(string)
	default:
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	threadID, err := strconv.Atoi(threadname)
	var thread *entity.Thread
	switch err {
	case nil:
		thread, err = threadInfo.threadApp.GetThreadByID(threadID)
	default:
		thread, err = threadInfo.threadApp.GetThreadByThreadname(threadname)
	}

	if err != nil {
		switch err {
		case entity.ThreadNotFoundError:
			responseBody, err := json.Marshal(entity.MessageOutput{"Can't find thread"})
			if err != nil {
				ctx.SetStatusCode(http.StatusInternalServerError)
				return
			}

			ctx.SetStatusCode(http.StatusNotFound)
			ctx.SetContentType("application/json")
			ctx.SetBody(responseBody)
			return
		default:
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}
	}

	responseBody, err := json.Marshal(thread)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(http.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(responseBody)
}

func (threadInfo *ThreadInfo) EditThread(ctx *fasthttp.RequestCtx) {
	threadnameInterface := ctx.UserValue("threadnameOrID")
	var threadname string

	switch threadnameInterface.(type) {
	case string:
		threadname = threadnameInterface.(string)
	default:
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	threadInput := new(entity.ThreadEditInput)
	err := json.Unmarshal(ctx.Request.Body(), threadInput)
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	threadID, err := strconv.Atoi(threadname)
	var thread *entity.Thread
	switch err {
	case nil:
		thread, err = threadInfo.threadApp.GetThreadByID(threadID)
	default:
		thread, err = threadInfo.threadApp.GetThreadByThreadname(threadname)
	}
	if err != nil {
		switch err {
		case entity.ThreadNotFoundError:
			responseBody, err := json.Marshal(entity.MessageOutput{"Can't find thread"})
			if err != nil {
				ctx.SetStatusCode(http.StatusInternalServerError)
				return
			}

			ctx.SetStatusCode(http.StatusNotFound)
			ctx.SetContentType("application/json")
			ctx.SetBody(responseBody)
			return
		default:
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}
	}

	//TODO: validate

	if *threadInput == (entity.ThreadEditInput{}) { // No need for editing
		responseBody, err := json.Marshal(thread)
		if err != nil {
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}

		ctx.SetStatusCode(http.StatusOK)
		ctx.SetContentType("application/json")
		ctx.SetBody(responseBody)
		return
	}

	if threadInput.Title != "" {
		thread.Title = threadInput.Title
	}
	if threadInput.Message != "" {
		thread.Message = threadInput.Message
	}

	err = threadInfo.threadApp.EditThread(thread)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	responseBody, err := json.Marshal(thread)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(http.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(responseBody)
}

func (threadInfo *ThreadInfo) GetThreadPosts(ctx *fasthttp.RequestCtx) {
	threadInput, err := entity.QueryToThreadGetPostsInput(ctx.QueryArgs())
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	threadnameInterface := ctx.UserValue("threadnameOrID")
	var threadname string
	switch threadnameInterface.(type) {
	case string:
		threadname = threadnameInterface.(string)
	default:
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	threadID, err := strconv.Atoi(threadname)
	var posts []*entity.Post
	startTime := time.Now()
	switch err {
	case nil:
		posts, err = threadInfo.threadApp.GetPostsByThreadID(threadID, threadInput.SortMode, threadInput.Limit, threadInput.StartAfter, threadInput.Desc)
	default:
		posts, err = threadInfo.threadApp.GetPostsByThreadname(threadname, threadInput.SortMode, threadInput.Limit, threadInput.StartAfter, threadInput.Desc)
	}

	if time.Since(startTime) > 100*time.Millisecond {
		fmt.Println("___________________")
		fmt.Println(time.Since(startTime))
		fmt.Println(threadInput)
		fmt.Println("___________________")
	}

	if err != nil && err != entity.PostNotFoundError {
		switch err {
		case entity.ThreadNotFoundError:
			responseBody, err := json.Marshal(entity.MessageOutput{"Can't find thread"})
			if err != nil {
				ctx.SetStatusCode(http.StatusInternalServerError)
				return
			}

			ctx.SetStatusCode(http.StatusNotFound)
			ctx.SetContentType("application/json")
			ctx.SetBody(responseBody)
			return
		case entity.UnsupportedSortingModeError:
			ctx.SetStatusCode(http.StatusBadRequest)
			return
		default:
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}
	}

	if posts == nil {
		posts = make([]*entity.Post, 0) // So that it marshalls as [] and not nil
	}

	responseBody, err := json.Marshal(posts)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(http.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(responseBody)
}

func (threadInfo *ThreadInfo) VoteThread(ctx *fasthttp.RequestCtx) {
	threadnameInterface := ctx.UserValue("threadnameOrID")
	var threadname string

	switch threadnameInterface.(type) {
	case string:
		threadname = threadnameInterface.(string)
	default:
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	voteInput := new(entity.VoteInput)
	err := json.Unmarshal(ctx.Request.Body(), voteInput)
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	threadID, err := strconv.Atoi(threadname)
	var thread *entity.Thread
	switch err {
	case nil:
		thread, err = threadInfo.threadApp.VoteThreadByID(threadID, voteInput.Username, voteInput.Vote)
	default:
		thread, err = threadInfo.threadApp.VoteThreadByThreadname(threadname, voteInput.Username, voteInput.Vote)
	}
	if err != nil {
		switch err {
		case entity.ThreadNotFoundError:
			responseBody, err := json.Marshal(entity.MessageOutput{"Can't find thread"})
			if err != nil {
				ctx.SetStatusCode(http.StatusInternalServerError)
				return
			}

			ctx.SetStatusCode(http.StatusNotFound)
			ctx.SetContentType("application/json")
			ctx.SetBody(responseBody)
			return
		case entity.UserNotFoundError:
			responseBody, err := json.Marshal(entity.MessageOutput{"Could not find user"})
			if err != nil {
				ctx.SetStatusCode(http.StatusInternalServerError)
				return
			}

			ctx.SetStatusCode(http.StatusNotFound)
			ctx.SetContentType("application/json")
			ctx.SetBody(responseBody)
			return
		case entity.IncorrectVoteAmountError:
			ctx.SetStatusCode(http.StatusBadRequest)
			return
		default:
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}
	}

	responseBody, err := json.Marshal(thread)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(http.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(responseBody)
}
