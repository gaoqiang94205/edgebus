package leader

import (
	"edgebus/common"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"log"
	"net"
	"net/http"
	"strings"
)

type Controller struct {
	Clients map[string]*websocket.Conn
}

//var Cs = make(map[string]*websocket.Conn)

// use default options
var upgrade = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
	return true
}}

func (l *Controller) Accept(w http.ResponseWriter, r *http.Request) {
	//升级为ws连接
	c, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	//将客户端IP注册
	host, _, err := net.SplitHostPort(c.RemoteAddr().String())
	//解析host异常进行处理
	if err != nil {
		log.Println("解析host异常")
		return
	}
	l.Clients[host] = c

	//服务端起单独协程进行消息的接受处理,完成后将连接关闭
	go func() {
		log.Println("receive data from client:" + host)
		receive(c)
	}()
}

//向边缘端发起请求
func (l *Controller) Deliver(w http.ResponseWriter, r *http.Request) {
	//获取Get请求查询参数
	query := r.URL.Query()
	//发送到指定的edge节点进行操作
	edge := query.Get(common.Edge)
	conn := l.Clients[edge]
	var build strings.Builder
	for k, _ := range l.Clients {
		build.WriteString(k)
	}
	if conn == nil {
		log.Printf("connection is invalid,check valid connection of:%s", build.String())
		return
	}
	//拼接一个请求
	request := common.SideRequest{
		Target: conn.RemoteAddr().String(),
		Url:    query.Get(common.Url),
		Method: query.Get(common.Method),
		Body:   "",
		Type:   common.Application,
	}
	data, err := json.Marshal(request)
	if err != nil {
		logrus.Errorf("marshal data with err: %v", err)
		return
	}
	if err := conn.WriteMessage(1, data); err != nil {
		logrus.Errorf("write data to client with err: %v", err)
	}
}

func (l *Controller) Ping(w http.ResponseWriter, r *http.Request) {
	logrus.Info("send a request from %s" + r.RemoteAddr)
	w.WriteHeader(200)
	w.Write([]byte("pong"))
}

//从边缘端接受消息并进行处理
func receive(c *websocket.Conn) {
	defer c.Close()
	for {
		//接收到请求
		_, message, err := c.ReadMessage()
		if err != nil {
			logrus.Errorf("receive data with err: %v", err)
		}
		//处理请求逻辑
		logrus.WithFields(logrus.Fields{"message": message}).Info("receive data from client")
	}
}
