package keepalive

import (
	"edgebus/client/pkg"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	heartbeatInterval = 25 * time.Second
	maxHeartbeat      = 3
)

//启动定时器进行心跳检测
func StartHeartbeat() {
	go func() {
		beat := 0
		ticker := time.NewTicker(heartbeatInterval)
		defer ticker.Stop()
		for {
			<-ticker.C
			//发送心跳
			agent := pkg.EdgeAgent
			if err := agent.Socket.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(time.Second)); err != nil {
				//Manager.DisConnect <- agent
				beat++
				logrus.Errorf("send heart message to cloud:%v agent with err:%v,the number of failure is %d",
					agent.Cloud, err, beat)
				if beat > maxHeartbeat {
					handleHeartFailure()
				}
			} else {
				beat = 0
			}
		}
	}()
}

//处理心率衰竭异常重新连接云端代理
func handleHeartFailure() {
	agent := pkg.EdgeAgent
	//将原agent置为nil
	pkg.EdgeAgent = nil
	if err := pkg.InitAgent(agent.Cloud.ToInfoRequest()); err != nil {
		logrus.Errorf("lose connection with ws server for err:%v,please check", err)
	}
}
