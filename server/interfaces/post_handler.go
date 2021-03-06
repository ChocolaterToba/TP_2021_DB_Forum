package interfaces

import (
	"dbforum/application"
	"dbforum/domain/entity"
	"fmt"
	"net/http"
	"strconv"

	json "github.com/mailru/easyjson"

	"github.com/valyala/fasthttp"
)

type PostInfo struct {
	postApp   application.PostAppInterface
	threadApp application.ThreadAppInterface
}

func NewPostInfo(postApp application.PostAppInterface, threadApp application.ThreadAppInterface) *PostInfo {
	return &PostInfo{
		postApp:   postApp,
		threadApp: threadApp,
	}
}

func (postInfo *PostInfo) CreatePost(ctx *fasthttp.RequestCtx) {
	threadnameInterface := ctx.UserValue("threadnameOrID")
	var threadname string

	switch threadnameInterface.(type) {
	case string:
		threadname = threadnameInterface.(string)
	default:
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	posts := entity.Posts(make([]*entity.Post, 0))

	err := json.Unmarshal(ctx.Request.Body(), &posts)
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	threadID, err := strconv.Atoi(threadname)
	var thread *entity.Thread
	switch err {
	case nil:
		thread, err = postInfo.threadApp.GetThreadByID(threadID)
	default:
		thread, err = postInfo.threadApp.GetThreadByThreadname(threadname)
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
	newPosts, err := postInfo.postApp.CreatePosts(posts, thread)
	if err != nil {
		switch err {
		case entity.ParentNotFoundError:
			// responseBody, err := json.Marshal(entity.MessageOutput{"Can't find parent post"})
			// if err != nil {
			// 	ctx.SetStatusCode(http.StatusInternalServerError)
			// 	return
			// }

			// ctx.SetStatusCode(http.StatusNotFound)
			// ctx.SetContentType("application/json")
			// ctx.SetBody(responseBody)
			// return
			fallthrough // API is weird there
		case entity.ParentInAnotherThreadError:
			responseBody, err := json.Marshal(entity.MessageOutput{"Parent post is in another thread"})
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
			fmt.Println(err)
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}
	}

	responseBody, err := json.Marshal(entity.Posts(newPosts))
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(http.StatusCreated)
	ctx.SetContentType("application/json")
	ctx.SetBody(responseBody)
}

func (postInfo *PostInfo) GetPost(ctx *fasthttp.RequestCtx) {
	postIDInterface := ctx.UserValue("postID")
	postID := 0

	var err error
	switch postIDInterface.(type) {
	case string:
		postID, err = strconv.Atoi(postIDInterface.(string))
		if err != nil {
			ctx.SetStatusCode(http.StatusBadRequest)
			return
		}
	default:
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	relatedObjects, err := entity.QueryToRelatedObjectsInput(ctx.QueryArgs())
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	post, err := postInfo.postApp.GetPostByID(postID)
	if err != nil {
		switch err {
		case entity.PostNotFoundError:
			responseBody, err := json.Marshal(entity.MessageOutput{"Can't find post"})
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
	postWithRelatives, err := postInfo.postApp.GetPostRelatives(post, relatedObjects.RelatedObjects)

	responseBody, err := json.Marshal(postWithRelatives)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(http.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(responseBody)
}

func (postInfo *PostInfo) EditPost(ctx *fasthttp.RequestCtx) {
	postIDInterface := ctx.UserValue("postID")
	postID := 0

	var err error
	switch postIDInterface.(type) {
	case string:
		postID, err = strconv.Atoi(postIDInterface.(string))
		if err != nil {
			ctx.SetStatusCode(http.StatusBadRequest)
			return
		}
	default:
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	postInput := new(entity.PostEditInput)
	err = json.Unmarshal(ctx.Request.Body(), postInput)
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	post, err := postInfo.postApp.GetPostByID(postID)
	if err != nil {
		switch err {
		case entity.PostNotFoundError:
			responseBody, err := json.Marshal(entity.MessageOutput{"Can't find post"})
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

	if postInput.Message == "" || post.Message == postInput.Message { // No need for editing
		responseBody, err := json.Marshal(post)
		if err != nil {
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}

		ctx.SetStatusCode(http.StatusOK)
		ctx.SetContentType("application/json")
		ctx.SetBody(responseBody)
		return
	}

	post.Message = postInput.Message
	post.IsEdited = true

	err = postInfo.postApp.EditPost(post)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	responseBody, err := json.Marshal(post)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(http.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(responseBody)
}
