package interfaces

import (
	"dbforum/application"
	"net/http"

	json "github.com/mailru/easyjson"

	"github.com/valyala/fasthttp"
)

type ServiceInfo struct {
	serviceApp application.ServiceAppInterface
}

func NewServiceInfo(serviceApp application.ServiceAppInterface) *ServiceInfo {
	return &ServiceInfo{serviceApp}
}

func (serviceInfo *ServiceInfo) GetForumStats(ctx *fasthttp.RequestCtx) {
	service, err := serviceInfo.serviceApp.GetStats()
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	responseBody, err := json.Marshal(service)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(http.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(responseBody)
}

func (serviceInfo *ServiceInfo) ClearForum(ctx *fasthttp.RequestCtx) {
	err := serviceInfo.serviceApp.TruncateAll()
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(http.StatusOK)
}
