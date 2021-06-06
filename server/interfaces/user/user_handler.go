package user

import (
	"dbforum/application"
	"dbforum/domain/entity"
	"encoding/json"
	"net/http"

	"github.com/valyala/fasthttp"
)

type UserInfo struct {
	userApp application.UserAppInterface
}

func NewUserInfo(userApp application.UserAppInterface) *UserInfo {
	return &UserInfo{userApp}
}

func (userInfo *UserInfo) CreateUser(ctx *fasthttp.RequestCtx) {
	usernameInterface := ctx.UserValue("username")
	userInput := new(entity.User)

	switch usernameInterface.(type) {
	case string:
		userInput.Username = usernameInterface.(string)
	default:
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	err := json.Unmarshal(ctx.Request.Body(), userInput)
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	//TODO: validate

	createUserInfo, err := userInfo.userApp.CreateUser(userInput)
	if err != nil {
		switch err {
		case entity.UserConflictError:
			responseBody, err := json.Marshal(createUserInfo.([]*entity.User))
			if err != nil {
				ctx.SetStatusCode(http.StatusInternalServerError)
				return
			}

			ctx.SetStatusCode(http.StatusConflict)
			ctx.SetBody(responseBody)
			return
		}

		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	responseBody, err := json.Marshal(userInput)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusCreated)
	ctx.SetBody(responseBody)
}

func (userInfo *UserInfo) GetUser(ctx *fasthttp.RequestCtx) {
	usernameInterface := ctx.UserValue("username")
	var username string

	switch usernameInterface.(type) {
	case string:
		username = ctx.UserValue("username").(string)
	default:
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	user, err := userInfo.userApp.GetUserByUsername(username)
	if err != nil {
		switch err {
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

	responseBody, err := json.Marshal(user)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(responseBody)
}

func (userInfo *UserInfo) EditUser(ctx *fasthttp.RequestCtx) {
	usernameInterface := ctx.UserValue("username")
	userInput := new(entity.User)

	switch usernameInterface.(type) {
	case string:
		userInput.Username = ctx.UserValue("username").(string)
	default:
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	err := json.Unmarshal(ctx.Request.Body(), userInput)
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	//TODO: validate

	err = userInfo.userApp.EditUser(userInput)
	if err != nil {
		switch err {
		case entity.UserNotFoundError:
			responseBody, err := json.Marshal(entity.MessageOutput{"Can't find user"})
			if err != nil {
				ctx.SetStatusCode(http.StatusInternalServerError)
				return
			}

			ctx.SetStatusCode(http.StatusNotFound)
			ctx.SetBody(responseBody)
			return

		case entity.UserConflictError:
			responseBody, err := json.Marshal(entity.MessageOutput{"Found conflicting user"})
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

	responseBody, err := json.Marshal(userInput)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(responseBody)
}
