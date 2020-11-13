package main

import (
	"edgebus/server/pkg"
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var addr = flag.String("addr", "localhost:8090", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)
	ld := pkg.NewLeader(make(map[string]*websocket.Conn))
	http.HandleFunc("/ws", ld.Accept)
	http.HandleFunc("/send",ld.Deliver)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
