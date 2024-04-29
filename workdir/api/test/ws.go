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

func testWs(ctx *gin.Context) {
	conn, err := wsupgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Printf("Failed to set websocket upgrade: %+v\n", err)
		return
	}

	defer conn.Close()

	for {
		// メッセージ読み込み
		mt, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("Error during message reading: %v\n", err)
			break
		}
		fmt.Printf("Received: %s\n", message)

		// エコーとしてメッセージを返す
		err = conn.WriteMessage(mt, message)
		if err != nil {
			fmt.Printf("Error during message writing: %v\n", err)
			break
		}
	}
}
