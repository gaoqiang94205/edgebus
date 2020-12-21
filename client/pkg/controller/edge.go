package controller

import (
	"edgebus/client/pkg"
	"edgebus/common"
	"github.com/gin-gonic/gin"
)

type Controller struct {
}

func (c *Controller) Connect(ctx *gin.Context) {
	cloud := &common.CloudInfoRequest{}
	ctx.BindJSON(cloud)
	pkg.InitAgent(*cloud)
}
