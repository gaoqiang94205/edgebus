package controller

import (
	"edgebus/common"
	"edgebus/server/pkg/connect"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"net/http"
)

type AcceptController struct {
}

func (l *AcceptController) Accept(ctx *gin.Context) {
	//升级http连接为websocket
	upGrader := websocket.Upgrader{
		// cross origin domain
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		// 处理 Sec-WebSocket-Protocol Header
		Subprotocols: []string{ctx.GetHeader("Sec-WebSocket-Protocol")},
	}
	conn, err := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		logrus.Printf("websocket connect error: %s", ctx.Param("channel"))
		return
	}
	ur := &common.UpgradeRequest{}
	ctx.BindJSON(ur)
	client := &connect.Client{
		Socket: conn,
		Name:   ur.Name,
		Id:     uuid.NewV4().String(),
	}
	connect.RegisterClient(client)
}

func (l *AcceptController) Test(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": gin.H{
			"store_result": "hello word",
		},
	})
}

//向边缘端发起请求
func (l *AcceptController) Deliver(c *gin.Context) {
	//获取Get请求查询参数
	//发送到指定的edge节点进行操作
	//edge := c.Query(common.Edge)

	//conn := Cs[edge].conn
	//var build strings.Builder
	//for k, _ := range Cs {
	//	build.WriteString(k)
	//}

	//拼接一个请求
	//request := common.SideRequest{
	//	//Target: conn.RemoteAddr().String(),
	//	Url:    c.Query(common.Url),
	//	Method: c.Query(common.Method),
	//	Body:   "",
	//	Type:   common.Application,
	//}
	//data, err := json.Marshal(request)
	//if err != nil {
	//	logrus.Errorf("marshal data with err: %v", err)
	//	return
	//}
	//if err := conn.WriteMessage(1, data); err != nil {
	//	logrus.Errorf("write data to client with err: %v", err)
	//}
}

func (l *AcceptController) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"error_code": 0,
		"message":    "ok",
		"data":       "pong",
	})

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
