package apiTest

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// クロスオリジン許可
	CheckOrigin: func(req *http.Request) bool {
		return true
	},
}

type Message struct {
	Type    int
	Message []byte
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)

func testWs(ctx *gin.Context) {
	conn, err := wsupgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Printf("Failed to set websocket upgrade: %+v\n", err)
		return
	}
	clients[conn] = true

	// 読み込み時エラーなどで処理が終わってもconnを閉じて削除
	defer conn.Close()
	defer delete(clients, conn)

	for {
		// メッセージ読み込み
		mt, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("Error during message reading: %v\n", err)
			break
		}
		broadcast <- Message{Type: mt, Message: message}
		fmt.Printf("Received: %s\n", message)

		// エコーとしてメッセージを返す
		// err = conn.WriteMessage(mt, message)
		// if err != nil {
		// 	fmt.Printf("Error during message writing: %v\n", err)
		// 	break
		// }
	}
}

func WsHandleMessages() {
	for {
		message := <-broadcast
		for client := range clients {
			err := client.WriteMessage(message.Type, message.Message)
			if err != nil {
				fmt.Printf("Error during message writing: %v\n", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
