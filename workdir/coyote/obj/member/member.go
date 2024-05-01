package memberObj

import "github.com/gorilla/websocket"

type Member struct {
	ID     int
	Name   string
	Client *websocket.Conn
}

var members = make(map[int]Member)
