package memberObj

import (
	"math/rand"

	"github.com/gorilla/websocket"
)

type Member struct {
	Name    string          `json:"name"`
	ImageID int             `json:"image_id"`
	Conn    *websocket.Conn `json:"-"`
}

// var members = make(map[string]*Member)

func CreateMember(conn *websocket.Conn, name string) Member {

	return Member{name, rand.Intn(67), conn}
}

// func GetMember(name string) *Member {
// 	return members[name]
// }

// func RemoveMember(name string) {
// 	delete(members, name)
// }
