package stateObj

import (
	cardObj "gin_test/coyote/obj/card.go"
	memberObj "gin_test/coyote/obj/member"
	roomObj "gin_test/coyote/obj/room"
	"gin_test/coyote/util"
)

type MemberStatus struct {
	ID     int              `json:"id"`
	Member memberObj.Member `json:"member"`
	Card   cardObj.Card     `json:"card"`
	Life   int              `json:"life"`
}

type State struct {
	StateNum              int                      `json:"state_num"` // 0:初期値()、1宣言フェーズ、2:計算フェーズ
	Table                 map[string]*MemberStatus `json:"table"`     // key: member_name
	DeclaredNum           int                      `json:"declared_num"`
	DeclaredMemberName    string                   `json:"declared_member_name"`
	NextDeclareMemberName string                   `json:"next_declare_member_name"`
	LimitNum              int                      `json:"limit_num"`
	TurnCount             int                      `json:"turn_count"`
	EndAccepts            []string                 `json:"end_accepts"`
}

// key: roomId
var states = make(map[string]*State)

func CreateState(roomId string, deck *cardObj.Deck) State {
	room := roomObj.GetRoomMemoryByID(roomId)
	table := createTable(*room, deck)
	state := State{0, table, 0, "", room.Members[0].Name, 0, 0, make([]string, 0)}
	states[roomId] = &state

	util.Log(util.LogObj{"log(create state)", state})
	return state
}
func GetStateFromMemory(roomId string) *State {
	if state, ok := states[roomId]; ok {
		util.Log(util.LogObj{"log(get state)", state})
		return state
	}
	return nil
}

func (state *State) RemoveMemberStatus(member memberObj.Member) {
	util.Log(util.LogObj{"log(remove member status from state)", map[string]interface{}{"state": state, "member": member}})
	delete(state.Table, member.Name)
}

func (state *State) ProceedStateToDeclare(roomId string, reqState State) {
	state.StateNum = 1
	state.DeclaredNum = reqState.DeclaredNum
	state.DeclaredMemberName = reqState.DeclaredMemberName
	state.NextDeclareMemberName = state.getNextMemberName(reqState.DeclaredMemberName)
}

func (state *State) ProceedStateToCalc(deck *cardObj.Deck, reqState State) {
	state.StateNum = 2
	// リミット計算
	state.LimitNum = state.calcLimit(deck)
}
func (state *State) DecrementMemberStatusLife(memberName string) {
	for name, memberStatus := range state.Table {
		if name == memberName {
			memberStatus.Life--
			if memberStatus.Life == 0 {
				state.NextDeclareMemberName = state.getNextMemberName(memberStatus.Member.Name)
			}
			return
		}
	}
}
func (state *State) GetSurvivers() []*MemberStatus {
	survivers := make([]*MemberStatus, 0)
	for _, memberStatus := range state.Table {
		if memberStatus.Life > 0 {
			survivers = append(survivers, memberStatus)
		}
	}
	return survivers
}

func (state *State) ProceedStateToInit(roomId string, deck *cardObj.Deck) {
	state.StateNum = 0
	state.DeclaredNum = 0
	state.DeclaredMemberName = state.NextDeclareMemberName
	state.NextDeclareMemberName = state.getNextMemberName(state.NextDeclareMemberName)
	state.TurnCount++
	state.EndAccepts = make([]string, 0)

	state.updateTable(deck)
}

func (state *State) AddEndAccepts(memberName string) {
	state.EndAccepts = append(state.EndAccepts, memberName)
}

/* local func */
func createTable(room roomObj.Room, deck *cardObj.Deck) map[string]*MemberStatus {
	table := make(map[string]*MemberStatus)
	for i, member := range room.Members {
		memberStatus := MemberStatus{i, member, deck.DrawCard(), 3}
		table[member.Name] = &memberStatus
	}
	util.Log(util.LogObj{"log(create table in state)", table})
	return table
}

func (state *State) updateTable(deck *cardObj.Deck) {
	// シャッフルフラグが立っていた場合シャッフル
	if deck.IsShuffle {
		deck.Shuffle()
	}

	for _, memberStatus := range state.Table {
		memberStatus.Card = deck.DrawCard()
	}
	util.Log(util.LogObj{"log(update table in state)", state.Table})
}

func (state *State) calcLimit(deck *cardObj.Deck) int {
	isDouble, isMaxZero := false, false
	// 特殊カード等を事前処理
	commonCards, specialCards := preprocessCalcLimit(state.Table, deck)

	for _, card := range specialCards {
		switch card.ID {
		case 200:
			isDouble = true
		case 201:
			isMaxZero = true
		}
	}

	var limit int
	var maxNum int = calcMaxNum(commonCards)
	for _, card := range commonCards {
		num := card.Num
		if isMaxZero && num == maxNum { // 最大値が0の場合は発火しない
			isMaxZero = false // 最大値１つのみを更新
			num = 0
		}
		if isDouble {
			num = num * 2
		}
		limit += num
	}
	util.Log(util.LogObj{"log(result to calcLimit)", limit})
	return limit
}

// ゲームフロー影響カードを発火、計算時影響カードを抽出
func preprocessCalcLimit(table map[string]*MemberStatus, deck *cardObj.Deck) ([]cardObj.Card, []cardObj.Card) {
	var common []cardObj.Card
	var special []cardObj.Card

	for name, memberStatus := range table {
		// ※洞穴カード:101 が出たら該当のmemberStatusを更新（引き直し）
		if memberStatus.Card.ID == 101 {
			table[name].Card = deck.DrawCard()
			util.Log(util.LogObj{"log(in preprocessCalcLimit(), cardId 101 fired)", table[name]})
		}

		card := table[name].Card
		if card.ID < 200 {
			common = append(common, card)
		} else {
			special = append(special, card)
		}
	}

	return common, special
}

func calcMaxNum(cards []cardObj.Card) int {
	var max int = -1000
	for _, card := range cards {
		if max < card.Num {
			max = card.Num
		}
	}
	return max
}

func (state *State) getNextMemberName(memberName string) string {
	currentMemberStatus := state.Table[memberName]
	var nextMemberStatus *MemberStatus

	// currentMemberStatusよりもIDが大きく、その中でIDが最小のものを次のメンバーへ
	for _, memberStatus := range state.Table {
		if memberStatus.Life == 0 {
			continue
		}
		if currentMemberStatus.ID < memberStatus.ID && nextMemberStatus == nil {
			nextMemberStatus = memberStatus
		}
		if nextMemberStatus != nil && currentMemberStatus.ID < memberStatus.ID && memberStatus.ID < nextMemberStatus.ID {
			nextMemberStatus = memberStatus
		}
	}

	// ↑の判定でもnextMemberStatusがnilの場合、IDが最小のものを次メンバーへ
	if nextMemberStatus == nil {
		for _, memberStatus := range state.Table {
			if memberStatus.Life == 0 {
				continue
			}
			nextMemberStatus = memberStatus
			break
		}
		for _, memberStatus := range state.Table {
			if memberStatus.Life == 0 {
				continue
			}
			if nextMemberStatus.ID > memberStatus.ID {
				nextMemberStatus = memberStatus
			}
		}
	}

	return nextMemberStatus.Member.Name
}
