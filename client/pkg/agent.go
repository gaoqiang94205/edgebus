package pkg

import (
	"edgebus/client/pkg/util"
	"edgebus/common"
	"encoding/json"
	"flag"
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
)

var addr = flag.String("addr", "localhost:8090", "http service address")

var Con *websocket.Conn

func Agent() error {
	//系统中断信号
	//interrupt := make(chan os.Signal, 1)
	//signal.Notify(interrupt, os.Interrupt)
	//建立连接
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	logrus.Info("connecting to %s", u.String())

	var err error
	Con, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}
	//receive执行完毕之后关闭该连接,理论上会一直保持连接
	defer Con.Close()

	receive()
	return nil
}

//循环处理来自peer对端来的请求，通过chan管道控制执行结束
func receive() {
	for {
		_, message, err := Con.ReadMessage()
		if err != nil {
			log.Println("read:", err)
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
		if err := Con.WriteMessage(1, body); err != nil {
			fmt.Errorf("write message fail")
		}
	}
}
func parseAppMessage(re *common.SideRequest) ([]byte, error) {
	client, err := util.GetURLClient(nil)
	//错误做记录然后抛出
	if err != nil {
		log.Println("create http client fail")
		return nil, err
	}
	resp, err := client.HTTPDo("GET", re.Url, http.Header{}, nil)
	if err != nil {
		log.Println("send http request with error" + err.Error())
		return nil, err
	}
	respBodyReader := http.MaxBytesReader(nil, resp.Body, maxBodySize)
	body, err := ioutil.ReadAll(respBodyReader)
	if err != nil {
		fmt.Errorf("read fail %s", err)
		return nil, err
	}
	return body, nil
}
