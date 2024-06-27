package sessionObj

import (
	cardObj "gin_test/coyote/obj/card.go"
	memberObj "gin_test/coyote/obj/member"
	stateObj "gin_test/coyote/obj/state"
	"gin_test/coyote/util"
)

type Session struct {
	EndAccepts []string        `json:"end_accepts"`
	State      *stateObj.State `json:"state"`
	Deck       cardObj.Deck    `json:"deck"`
}

// key: roomId
var sessions map[string]*Session = make(map[string]*Session)

func CreateSession(roomId string) Session {
	deck := cardObj.CreateDeck()
	state := stateObj.CreateState(roomId, &deck)
	session := Session{make([]string, 0), &state, deck}
	sessions[roomId] = &session

	util.Log(util.LogObj{"log(create session)", session})
	return session
}

func GetSessionFromMemory(roomId string) *Session {
	if session, ok := sessions[roomId]; ok {
		util.Log(util.LogObj{"log(get session)", session})
		return session
	}
	return nil
}

func RemoveStateMember(roomId string, member memberObj.Member) {

}
