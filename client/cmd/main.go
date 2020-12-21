package main

import (
	"edgebus/client/pkg"
	"edgebus/client/pkg/routes"
	"edgebus/common"
	"edgebus/common/http"
	st "edgebus/pkg/setting"
	"edgebus/tools/log"
	"flag"
	"github.com/sirupsen/logrus"
	"os"
)

const (
	ClientType = "client"
	ServerType = "server"
)

func init() {
	log.Setup()
	st.SetClient()
}

func main() {
	flag.Parse()
	cf := st.ClientConf
	//启动http服务
	router := http.InitRoutes()
	routes.RegisterApiRouter(router)
	router.Run(st.ServerConf.Listen())

	if cf.Type == ClientType && cf.Target != "" {
		pkg.InitAgent(common.CloudInfoRequest{Addr: cf.Target, Path: pkg.WsPath})
	} else if cf.Type == ServerType {

	} else {
		logrus.StandardLogger().Errorf("please check the args")
		os.Exit(1)
	}

}
