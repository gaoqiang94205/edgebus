package routes

import (
	"edgebus/server/pkg/controller"
	"github.com/gin-gonic/gin"
)

func RegisterApiRouter(router *gin.Engine) {
	ac := &controller.AcceptController{}
	apiRouter := router.Group("api/vcloud/v1/edge")
	{
		apiRouter.GET("/test", ac.Test)
		apiRouter.POST("/ws", ac.Accept)
	}
}
