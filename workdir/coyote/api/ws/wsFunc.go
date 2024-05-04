package coyoteWsApi

import (
	"encoding/json"
	memberObj "gin_test/coyote/obj/member"
	roomObj "gin_test/coyote/obj/room"
)

/*
type -1
Error
to 1
*/
func errorOccurred(errMsg string, member memberObj.Member, roomId string) {
	type Error struct {
		Message string `json:"message"`
	}
	type Data struct {
		Error Error `json:"error"`
	}
	marshaledData, _ := json.Marshal(Data{Error{errMsg}})

	member.Conn.WriteJSON(WSMessage{-1, "error", roomId, marshaledData})
}

/*
type 0
Members Update
to 多
*/
func membersUpdate(roomId string, broadcast chan WSMessage) {
	room := roomObj.GetRoomMemoryByID(roomId)
	marshaledData, _ := json.Marshal(map[string]interface{}{"room": room})
	broadcast <- WSMessage{0, "members update", roomId, marshaledData}
}

/*
type 1
Send Comment
to 多
*/
func sendComment(reqMsg WSMessage, roomId string, broadcast chan WSMessage) {
	// type Comment struct {
	//	Name string `json:"name"`
	// 	Content string `json:"content"`
	// }
	// type Data struct {
	// 	Comment Comment `json:"comment"`
	// }

	broadcast <- WSMessage{1, "send comment", roomId, reqMsg.Data}
}
