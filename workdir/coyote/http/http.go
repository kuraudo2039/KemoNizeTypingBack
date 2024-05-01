package coyoteHttpApi

import (
	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

func RegistHttpApi(engine *gin.Engine, client *firestore.Client) {
	engine.POST("/coyote/entry", enterRoom(client))
}
