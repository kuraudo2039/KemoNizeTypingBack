package roomObj

import (
	"context"
	memberObj "gin_test/coyote/obj/member"
	"gin_test/coyote/util"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type RoomData struct {
	Password string `firestore:"password"`
	State    int    `firestore:"state"`
}

type Room struct {
	ID      string              `json:id`
	Data    RoomData            `json:data`
	Members []*memberObj.Member `json:members`
	State   int                 `json:state`
}

var rooms = make(map[string]Room)

func CreateRoom(client *firestore.Client, ctx context.Context, data RoomData) (*Room, error) {
	resDoc, _, err := client.Collection("room").Add(ctx, data)
	util.Log(util.LogObj{"created room", resDoc})

	room := Room{resDoc.ID, data, make([]*memberObj.Member, 0), 0}
	rooms[resDoc.ID] = room
	return &room, err
}

func GetRoom(client *firestore.Client, ctx context.Context, pw string) (*Room, error) {
	query := client.Collection("room").Where("password", "==", pw)
	iter := query.Documents(ctx)

	doc, err := iter.Next()
	if err == iterator.Done {
		err = nil
	}

	var roomData RoomData
	// 一つでも該当roomがあれば、roomを返す
	if doc != nil {
		doc.DataTo(&roomData)
		var room Room
		if val, ok := rooms[doc.Ref.ID]; ok {
			room = val
		} else {
			room = Room{doc.Ref.ID, roomData, make([]*memberObj.Member, 0), 0}
		}

		util.Log(util.LogObj{"get room", room})
		return &room, nil
	} else {
		// 無ければ
		return nil, err
	}
}
