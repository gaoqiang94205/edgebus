package main

import (
	"edgebus/common/http"
	"edgebus/pkg/setting"
	"edgebus/server/pkg/routes"
	"edgebus/tools/log"
)

func init() {
	log.Setup()
	setting.SetServer()
}
func main() {
	//初始化路由信息
	router := http.InitRoutes()
	routes.RegisterApiRouter(router)
	//启动web服务
	router.Run(setting.ServerConf.Listen())
}
