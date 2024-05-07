package coyoteWsApi

import (
	"encoding/json"
	memberObj "gin_test/coyote/obj/member"
	roomObj "gin_test/coyote/obj/room"
	sessionObj "gin_test/coyote/obj/session"
	stateObj "gin_test/coyote/obj/state"
	"time"
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
	data := map[string]interface{}{"room": room}
	if session := sessionObj.GetSessionFromMemory(roomId); session != nil {
		data["session"] = session
	}
	marshaledData, _ := json.Marshal(data)
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

/*
type 10
Start Game
to 多
*/
func startSession(roomId string, broadcast chan WSMessage) {
	session := sessionObj.CreateSession(roomId)
	marshaledData, _ := json.Marshal(map[string]interface{}{"session": session})
	broadcast <- WSMessage{10, "start session", roomId, marshaledData}
}

/*
type 11
Declare Num
to 多
*/
func declareNum(reqMsg WSMessage, roomId string, broadcast chan WSMessage) {
	session := sessionObj.GetSessionFromMemory(roomId)
	type Data struct {
		State stateObj.State `json:"state"`
	}

	var data Data
	json.Unmarshal(reqMsg.Data, &data)
	session.State.ProceedStateToDeclare(roomId, data.State)

	marshaledData, _ := json.Marshal(map[string]interface{}{"session": session})
	broadcast <- WSMessage{11, "declare num", roomId, marshaledData}
}

/*
type 12
Declare Coyote
to 多
*/
func declareCoyote(reqMsg WSMessage, roomId string, broadcast chan WSMessage) {
	session := sessionObj.GetSessionFromMemory(roomId)
	type Data struct {
		State stateObj.State `json:"state"`
	}

	var data Data
	json.Unmarshal(reqMsg.Data, &data)
	// 1. 計算結果をステートに反映
	session.State.ProceedStateToCalc(&session.Deck, data.State)
	// 2. 計算結果をもとにコヨーテ成立判定
	if session.State.LimitNum < session.State.DeclaredNum {
		session.State.DecrementMemberStatusLife(session.State.DeclaredMemberName)
	} else {
		session.State.DecrementMemberStatusLife(session.State.NextDeclareMemberName)
	}
	marshaledData, _ := json.Marshal(map[string]interface{}{"session": session})

	// 3. 勝敗が決していたらtype:20として送信
	surviversCount := len(session.State.GetSurvivers())
	if surviversCount > 1 {
		broadcast <- WSMessage{12, "declare coyote", roomId, marshaledData}
		// 自動次ターン移行プロセス起動
		go autoProceedNextTurn(roomId, broadcast)
	} else {
		broadcast <- WSMessage{20, "session end", roomId, marshaledData}
	}
}

/*
type 13 (no receive)
proceed next turn
to 多
*/
func proceedNextTurn(session *sessionObj.Session, roomId string, broadcast chan WSMessage) {
	session.State.ProceedStateToInit(roomId, &session.Deck)
	marshaledJson, _ := json.Marshal(map[string]interface{}{"session": session})
	broadcast <- WSMessage{13, "proceed next turn", roomId, marshaledJson}
}

/*
type 13 (write only)
auto proceed next turn
20秒経過で自動で次ターンへ移行
*/
func autoProceedNextTurn(roomId string, broadcast chan WSMessage) {
	session := sessionObj.GetSessionFromMemory(roomId)
	time.Sleep(20 * time.Second)
	if len(session.State.EndAccepts) != len(session.State.Table) {
		proceedNextTurn(session, roomId, broadcast)
	}
}

/*
type 13
acceptStateEnd
賛成票が集まったら次ターンへ移行
*/
func acceptStateEnd(member memberObj.Member, roomId string, broadcast chan WSMessage) {
	session := sessionObj.GetSessionFromMemory(roomId)
	session.State.AddEndAccepts(member.Name)
	if len(session.State.EndAccepts) == len(session.State.Table) {
		proceedNextTurn(session, roomId, broadcast)
	}
}
