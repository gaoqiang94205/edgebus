package main

import (
	"edgebus/client/pkg"
	"edgebus/pkg/setting"
	"edgebus/tools/log"
	"fmt"
)

func init() {
	log.Setup()
	setting.SetClient()
}

func main() {
	if err := pkg.Agent(); err != nil {
		fmt.Errorf("start client with err:%v", err)
	}
}
