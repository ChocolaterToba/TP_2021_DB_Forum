package forum

import (
	"dbforum/application"
	"dbforum/domain/entity"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/valyala/fasthttp"
)

type ForumInfo struct {
	forumApp application.ForumAppInterface
}

func NewForumInfo(forumApp application.ForumAppInterface) *ForumInfo {
	return &ForumInfo{forumApp}
}

func (forumInfo *ForumInfo) CreateForum(ctx *fasthttp.RequestCtx) {
	forumInput := new(entity.Forum)
	err := json.Unmarshal(ctx.Request.Body(), forumInput)
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	//TODO: validate

	createForumInfo, err := forumInfo.forumApp.CreateForum(forumInput)
	if err != nil {
		fmt.Println(err)
		switch err {
		case entity.ForumConflictError:
			responseBody, err := json.Marshal(createForumInfo.(*entity.Forum))
			if err != nil {
				ctx.SetStatusCode(http.StatusInternalServerError)
				return
			}

			ctx.SetStatusCode(http.StatusConflict)
			ctx.SetBody(responseBody)
			return
		case entity.UserNotFoundError:
			responseBody, err := json.Marshal(entity.MessageOutput{"Can't find user"})
			if err != nil {
				ctx.SetStatusCode(http.StatusInternalServerError)
				return
			}

			ctx.SetStatusCode(http.StatusNotFound)
			ctx.SetBody(responseBody)
			return
		default:
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}
	}

	responseBody, err := json.Marshal(forumInput)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusCreated)
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

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(responseBody)
}
