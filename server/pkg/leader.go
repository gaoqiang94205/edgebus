package pkg

import (
	"edgebus/common"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/http"
)

type Leader struct {
	clients map[string]*websocket.Conn
}

func NewLeader(clients map[string]*websocket.Conn) *Leader {
	return &Leader{clients: clients}
}

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
	return true
}} // use default options

func (l *Leader) Accept(w http.ResponseWriter, r *http.Request) {
	//升级为ws连接
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	//将客户端IP注册
	host, _, err := net.SplitHostPort(c.RemoteAddr().String())
	//解析host异常进行处理
	if err!=nil{
		log.Println("解析host异常")
		return
	}

	l.clients[host] = c
	//服务端起单独协程进行消息的接受处理,完成后将连接关闭
	go func() {
		log.Println("receive data from client:" + host)
		receive(c)
	}()
}

//向边缘端发起请求
func (l *Leader) Deliver(w http.ResponseWriter, r *http.Request) {
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	conn := l.clients[host]
	if conn == nil {
		log.Println("connection is invalid")
		return
	}
	//拼接一个请求
	query := r.URL.Query()
	request := common.SideRequest{
		Target: conn.RemoteAddr().String(),
		Url:    query.Get(common.Url),
		Method: query.Get(common.Method),
		Body:   "",
		Type:   common.Application,
	}
	data, err := json.Marshal(request)
	if err != nil {
		log.Println("marshal data failed")
		return
	}
	if err := conn.WriteMessage(1, data); err != nil {
		log.Println("marshal data failed")
	}
}

//从边缘端接受消息并进行处理
func receive(c *websocket.Conn) {
	defer c.Close()
	for {
		//接收到请求
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read failed reason:", err)
			break
		}
		//处理请求逻辑
		log.Printf("recv: %s", string(message))
	}
}
