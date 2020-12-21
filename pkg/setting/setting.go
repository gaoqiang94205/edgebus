package setting

import (
	"flag"
	"fmt"
	"github.com/go-ini/ini"
	"log"
	"strings"
	"sync"
)

type global struct {
	LocalHost      string //本机内网IP
	ServerList     map[string]string
	ServerListLock sync.RWMutex
}

var GlobalSetting = &global{}

type clientConf struct {
	Type   string
	Target string
}

//func (c *clientConf) Call() string {
//	return c.Target
//}

type serverConf struct {
	ListenAddr string
	Port       string
}

func (s *serverConf) Listen() string {
	return strings.Join([]string{s.ListenAddr, s.Port}, ":")
}

var ClientConf = &clientConf{}
var ServerConf = &serverConf{}

var cfg *ini.File

func SetClient() {
	ClientConf.Type = *flag.String("type", "client", "--type client/server")
	ClientConf.Target = *flag.String("target", "127.0.0.1:8090", "127.0.0.1:8090")
}

func SetServer() {
	configFile := flag.String("c", "conf/ws-server.ini", "-c conf/server.ini")
	setup(configFile, ServerConf)
}

func setup(cff *string, cfg interface{}) {
	var err error
	cfg, err = ini.Load(cff)
	//加载配置文件失败直接异常退出os
	if err != nil {
		fmt.Errorf("parse ws client config fail")
	}
	ini.MapTo("client", cfg)
}

// mapTo map section
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}
