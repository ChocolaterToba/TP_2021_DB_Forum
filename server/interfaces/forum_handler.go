package interfaces

import (
	"dbforum/application"
	"dbforum/domain/entity"
	"net/http"

	json "github.com/mailru/easyjson"

	"github.com/valyala/fasthttp"
)

type ForumInfo struct {
	forumApp application.ForumAppInterface
}

func NewForumInfo(forumApp application.ForumAppInterface) *ForumInfo {
	return &ForumInfo{forumApp}
}

func (forumInfo *ForumInfo) CreateForum(ctx *fasthttp.RequestCtx) {
	forumInput := new(entity.ForumCreateInput)
	err := json.Unmarshal(ctx.Request.Body(), forumInput)
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	//TODO: validate

	createdForum, err := forumInfo.forumApp.CreateForum(forumInput)
	if err != nil {
		switch err {
		case entity.ForumConflictError:
			responseBody, err := json.Marshal(createdForum)
			if err != nil {
				ctx.SetStatusCode(http.StatusInternalServerError)
				return
			}

			ctx.SetStatusCode(http.StatusConflict)
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

	responseBody, err := json.Marshal(createdForum)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(http.StatusCreated)
	ctx.SetContentType("application/json")
	ctx.SetBody(responseBody)
}

func (forumInfo *ForumInfo) GetForum(ctx *fasthttp.RequestCtx) {
	forumnameInterface := ctx.UserValue("forumname")
	forumInput := new(entity.Forum)

	switch forumnameInterface.(type) {
	case string:
		forumInput.Forumname = forumnameInterface.(string)
	default:
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	forum, err := forumInfo.forumApp.GetForumByForumname(forumInput.Forumname)
	if err != nil {
		switch err {
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
		}

		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	responseBody, err := json.Marshal(forum)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(http.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(responseBody)
}

func (forumInfo *ForumInfo) GetForumUsers(ctx *fasthttp.RequestCtx) {
	forumnameInterface := ctx.UserValue("forumname")

	var forumname string
	switch forumnameInterface.(type) {
	case string:
		forumname = forumnameInterface.(string)
	default:
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	forumInput, err := entity.QueryToForumGetUsersInput(ctx.QueryArgs())
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	users, err := forumInfo.forumApp.GetUsersByForumname(forumname, forumInput.Limit, forumInput.StartAfter, forumInput.Desc)
	if err != nil && err != entity.UserNotFoundError {
		switch err {
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
		}

		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	if users == nil {
		users = make([]*entity.User, 0) // So that it marshalls as [] and not nil
	}

	responseBody, err := json.Marshal(entity.Users(users))
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(http.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(responseBody)
}

func (forumInfo *ForumInfo) GetForumThreads(ctx *fasthttp.RequestCtx) {
	forumnameInterface := ctx.UserValue("forumname")

	var forumname string
	switch forumnameInterface.(type) {
	case string:
		forumname = forumnameInterface.(string)
	default:
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	forumInput, err := entity.QueryToForumGetThreadsInput(ctx.QueryArgs())
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	threads, err := forumInfo.forumApp.GetThreadsByForumname(forumname, forumInput.Limit, forumInput.StartFrom, forumInput.Desc)
	if err != nil && err != entity.ThreadNotFoundError {
		switch err {
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
		}

		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	if threads == nil {
		threads = make([]*entity.Thread, 0) // So that it marshalls as [] and not nil
	}

	responseBody, err := json.Marshal(entity.Threads(threads))
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(http.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(responseBody)
}
