package memberObj

import "github.com/gorilla/websocket"

type Member struct {
	Name string          `json:"name"`
	Conn *websocket.Conn `json:"-"`
}

// var members = make(map[string]*Member)

func CreateMember(conn *websocket.Conn, name string) Member {

	return Member{name, conn}
}

// func GetMember(name string) *Member {
// 	return members[name]
// }

// func RemoveMember(name string) {
// 	delete(members, name)
// }
