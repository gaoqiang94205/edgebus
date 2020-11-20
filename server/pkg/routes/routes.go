package routes

import (
	"edgebus/server/pkg/leader"
	"github.com/gorilla/websocket"
	"net/http"
)

func Init() {
	ld := leader.Controller{make(map[string]*websocket.Conn)}

	http.HandleFunc("/ws", ld.Accept)
	http.HandleFunc("/send", ld.Deliver)
	http.HandleFunc("/ping", ld.Pong)
}
