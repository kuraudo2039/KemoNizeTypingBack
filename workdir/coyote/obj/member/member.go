package memberObj

import "github.com/gorilla/websocket"

type Member struct {
	Name string          `json:"name"`
	Conn *websocket.Conn `json:"-"`
}

// var members = make(map[int]Member)

func CreateMember(conn *websocket.Conn, name string) Member {
	return Member{name, conn}
}
