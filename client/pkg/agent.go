package pkg

import (
	"edgebus/client/pkg/util"
	"edgebus/common"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const (
	maxBodySize = 5 * 1e6
)

var addr = flag.String("addr", "localhost:8090", "http service address")

func Agent() {
	flag.Parse()
	log.SetFlags(0)

	//系统中断信号
	//interrupt := make(chan os.Signal, 1)
	//signal.Notify(interrupt, os.Interrupt)
	//建立连接
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	//defer c.Close()

	done := make(chan struct{})
	receive(done, c)
}

//循环处理来自peer对端来的请求，通过chan管道控制执行结束
func receive(signal chan struct{}, c *websocket.Conn) {
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
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
		if err := c.WriteMessage(1, body); err != nil {
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
