package pkg

import (
	"edgebus/client/pkg/keepalive"
	"edgebus/client/pkg/util"
	"edgebus/common"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const (
	maxBodySize = 5 * 1e6
	WsPath      = "/ws"
)

type agent struct {
	Socket *websocket.Conn
	Cloud  CloudInfo
}

type CloudInfo struct {
	Addr string
	Path string
}

func (cloud *CloudInfo) ToInfoRequest() common.CloudInfoRequest {
	return common.CloudInfoRequest{Addr: cloud.Addr, Path: cloud.Path}
}

var EdgeAgent *agent

func InitAgent(request common.CloudInfoRequest) error {
	//建立连接
	u := url.URL{Scheme: "ws", Host: request.Addr, Path: request.Path}
	logrus.Info("connecting to %s", u.String())

	con, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}
	EdgeAgent = &agent{Socket: con, Cloud: CloudInfo{Addr: request.Addr, Path: WsPath}}
	keepalive.StartHeartbeat()
	go receive()
	return nil
}

//循环处理来自peer对端来的请求，通过chan管道控制执行结束
func receive() {
	for {
		_, message, err := EdgeAgent.Socket.ReadMessage()
		if err != nil {
			logrus.Errorf("read data from cloud with err:%v", err)
			continue
		}
		re := new(common.SideRequest)
		if err := json.Unmarshal(message, re); err != nil {
			log.Println("json unmarshal failed")
		}
		var body []byte
		if re.Type == common.Application {
			body, err = parseAppMessage(re)
			if err != nil {
				log.Println("parse message failed")
				continue
			}
		} else if re.Type == common.Ping {
			body = []byte("pong")
		} else {
			body = []byte("hello")
		}
		if err := EdgeAgent.Socket.WriteMessage(1, body); err != nil {
			fmt.Errorf("write message fail")
		}
	}
}
func parseAppMessage(re *common.SideRequest) ([]byte, error) {
	client, err := util.GetURLClient(nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.HTTPDo("GET", re.Url, http.Header{}, nil)
	if err != nil {
		logrus.Errorf("send http request with error:%v", err)
		return nil, err
	}
	respBodyReader := http.MaxBytesReader(nil, resp.Body, maxBodySize)
	body, err := ioutil.ReadAll(respBodyReader)
	if err != nil {
		return nil, err
	}
	return body, nil
}
