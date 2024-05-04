package coyoteWsApi

import (
	"context"
	"encoding/json"
	memberObj "gin_test/coyote/obj/member"
	roomObj "gin_test/coyote/obj/room"
	"gin_test/coyote/util"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WSMessage struct {
	Type     int             `json:"type"`
	TypeName string          `json:"typename"`
	RoomID   string          `json:"roomid"`
	Data     json.RawMessage `json:"data"`
}

var broadcast = make(chan WSMessage)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// クロスオリジン許可
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ConnectWs(client *firestore.Client) func(*gin.Context) {
	return func(c *gin.Context) {
		// クエリパラメータ取得&ルーム存在チェック
		roomId := c.Query("id")
		memberName := c.Query("name")
		if room := roomObj.GetRoomMemoryByID(roomId); room == nil {
			util.Log(util.LogObj{"error(Failed to get room by roomId)", roomId})
			return
		}

		// websocketアップグレード
		conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			util.Log(util.LogObj{"error(Failed to websocket upgrade.)", err.Error()})
			return
		}
		defer conn.Close()

		// コンテキスト初期化
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// アイドルタイムアウトを30分に設定。アクティビティ時間と現在時間の差により切断
		idleTimeout := 30 * time.Minute
		lastActivity := time.Now()
		go checkIdleTimeout(conn, ctx, idleTimeout, &lastActivity)

		// メンバー作成＆ルームメンバー追加
		member := memberObj.CreateMember(conn, memberName)
		if room := roomObj.AddMember(roomId, member); room != nil {
			membersUpdate(roomId, broadcast)
		} else {
			util.Log(util.LogObj{"error(Failed to add member)", member})
			// errorOccurred("入室に失敗しました。\nルームが存在しないか、既に使われている名前です。", member, roomId)
			return
		}
		defer deferConnectWs(roomId, member)

		// イベントハンドリング開始
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				util.Log(util.LogObj{"error(Failed to read message)", err.Error()})
				break
			}
			lastActivity = time.Now() // 受信時にアクティビティ時間更新

			var msgJson WSMessage
			if err := json.Unmarshal(msg, &msgJson); err != nil {
				util.Log(util.LogObj{"error(Failed to unmarshal message)", err.Error()})
				continue
			}
			util.Log(util.LogObj{"received", msgJson})

			switch msgJson.Type {
			case 1:
				sendComment(msgJson, roomId, broadcast)
			case 2:

			default:
				util.Log(util.LogObj{"log(Unknown message type)", msg})
			}
		}
	}
}

func deferConnectWs(roomId string, member memberObj.Member) {
	util.Log(util.LogObj{"log", "launch deferConnectWs"})

	roomObj.RemoveMember(roomId, member)
	membersUpdate(roomId, broadcast)
}

func checkIdleTimeout(conn *websocket.Conn, ctx context.Context, idleTimeout time.Duration, lastActivity *time.Time) {
	for {
		time.Sleep(1 * time.Minute)
		select {
		case <-ctx.Done():
			util.Log(util.LogObj{"log(checkIdleTimeout routine shutdown)", lastActivity})
			return
		default:
			if time.Since(*lastActivity) > idleTimeout {
				conn.Close()
				return
			}
		}
	}
}

func checkHandleMessagesTimeout(cancel context.CancelFunc, idleTimeout time.Duration, lastActivity *time.Time) {
	for {
		time.Sleep(1 * time.Minute)
		if time.Since(*lastActivity) > idleTimeout {
			cancel()
			return
		}
	}
}

func HandleMessages() {
	ctx, cancel := context.WithCancel(context.Background())

	idleTimeout := 30 * time.Minute
	lastActivity := time.Now()

	go checkHandleMessagesTimeout(cancel, idleTimeout, &lastActivity)

	for {
		select {
		case <-ctx.Done():
			util.Log(util.LogObj{"log", "HandleMessages routine shutdown"})
			return
		case msg := <-broadcast:
			room := roomObj.GetRoomMemoryByID(msg.RoomID)
			for _, member := range room.Members {
				if err := member.Conn.WriteJSON(msg); err != nil {
					util.Log(util.LogObj{"error(Failed to write json at response)", err.Error()})
					roomObj.RemoveMember(msg.RoomID, member)
					continue
				}
			}
		}
	}
}
