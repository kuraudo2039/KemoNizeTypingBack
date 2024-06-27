package coyoteHttpApi

import (
	coyoteWsApi "gin_test/coyote/api/ws"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

func RegistHttpApi(engine *gin.Engine, client *firestore.Client) {
	engine.POST("/coyote/entry", enterRoom(client))
	engine.GET("/coyote/ws", coyoteWsApi.ConnectWs(client))
}
