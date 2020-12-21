package routes

import (
	"edgebus/client/pkg/controller"
	"github.com/gin-gonic/gin"
)

func RegisterApiRouter(router *gin.Engine) {
	c := &controller.Controller{}
	apiRouter := router.Group("api/vcloud/v1/edge")
	{
		apiRouter.GET("/test", c.Connect)
	}
}
