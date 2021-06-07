package interfaces

import (
	"dbforum/application"
	"dbforum/domain/entity"
	"encoding/json"
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
		threadInput.Created = time.Now()
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
	sortingMode := string(ctx.QueryArgs().Peek("sort"))
	switch sortingMode {
	case "flat", "tree", "parent_tree":
		// Intentional no-op
	case "":
		sortingMode = "flat"
	default:
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	sortingAscString := string(ctx.QueryArgs().Peek("desc"))
	sortingAsc := true
	switch sortingAscString {
	case "true":
		sortingAsc = false
	case "", "false":
		// Intentional no-op
	default:
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
	var posts interface{}
	switch err {
	case nil:
		posts, err = threadInfo.threadApp.GetPostsByThreadID(threadID, sortingMode, sortingAsc)
	default:
		posts, err = threadInfo.threadApp.GetPostsByThreadname(threadname, sortingMode, sortingAsc)
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
