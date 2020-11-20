package main

import (
	"edgebus/pkg/setting"
	"edgebus/server/pkg/routes"
	"edgebus/tools/log"
	"github.com/sirupsen/logrus"
	"net/http"
)

func init() {
	log.Setup()
	setting.SetServer()
}
func main() {
	//初始化路由信息
	routes.Init()

	//启动web服务
	logrus.Fatal(http.ListenAndServe(setting.ServerConf.Listen(), nil))
}
