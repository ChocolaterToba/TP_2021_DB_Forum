package interfaces

import (
	"dbforum/application"
	"dbforum/domain/entity"
	"net/http"

	json "github.com/mailru/easyjson"

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
			users := createUserInfo.([]*entity.User)
			responseBody, err := json.Marshal(entity.Users(users))
			if err != nil {
				ctx.SetStatusCode(http.StatusInternalServerError)
				return
			}

			ctx.SetStatusCode(http.StatusConflict)
			ctx.SetContentType("application/json")
			ctx.SetBody(responseBody)
			return
		default:
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}
	}

	responseBody, err := json.Marshal(userInput)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(http.StatusCreated)
	ctx.SetContentType("application/json")
	ctx.SetBody(responseBody)
}

func (userInfo *UserInfo) GetUser(ctx *fasthttp.RequestCtx) {
	usernameInterface := ctx.UserValue("username")
	var username string

	switch usernameInterface.(type) {
	case string:
		username = usernameInterface.(string)
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
			ctx.SetContentType("application/json")
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

	ctx.SetStatusCode(http.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(responseBody)
}

func (userInfo *UserInfo) EditUser(ctx *fasthttp.RequestCtx) {
	usernameInterface := ctx.UserValue("username")

	var username string
	switch usernameInterface.(type) {
	case string:
		username = usernameInterface.(string)
	default:
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	userInput := new(entity.UserEditInput)
	err := json.Unmarshal(ctx.Request.Body(), userInput)
	if err != nil {
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
			ctx.SetContentType("application/json")
			ctx.SetBody(responseBody)
			return
		default:
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}
	}

	if *userInput == (entity.UserEditInput{}) { // No need for editing
		responseBody, err := json.Marshal(user)
		if err != nil {
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}

		ctx.SetStatusCode(http.StatusOK)
		ctx.SetContentType("application/json")
		ctx.SetBody(responseBody)
		return
	}

	//TODO: validate
	if userInput.FullName != "" {
		user.FullName = userInput.FullName
	}
	if userInput.Description != "" {
		user.Description = userInput.Description
	}
	if userInput.EMail != "" {
		user.EMail = userInput.EMail
	}

	err = userInfo.userApp.EditUser(user)
	if err != nil {
		switch err {
		case entity.UserNotFoundError:
			responseBody, err := json.Marshal(entity.MessageOutput{"Can't find user"})
			if err != nil {
				ctx.SetStatusCode(http.StatusInternalServerError)
				return
			}

			ctx.SetStatusCode(http.StatusNotFound)
			ctx.SetContentType("application/json")
			ctx.SetBody(responseBody)
			return

		case entity.UserConflictError:
			responseBody, err := json.Marshal(entity.MessageOutput{"Found conflicting user"})
			if err != nil {
				ctx.SetStatusCode(http.StatusInternalServerError)
				return
			}

			ctx.SetStatusCode(http.StatusConflict)
			ctx.SetContentType("application/json")
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

	ctx.SetStatusCode(http.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(responseBody)
}
