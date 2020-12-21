package connect

import "github.com/gorilla/websocket"

type Client struct {
	Name   string
	Id     string
	Socket *websocket.Conn
}

func (client *Client) getId() string {
	return client.Name + client.Id
}

type WsManager struct {
	clients map[string]*Client
}
