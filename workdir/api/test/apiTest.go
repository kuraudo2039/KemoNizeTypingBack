package apiTest

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

// regist API
func RegistApi(engine *gin.Engine, client *firestore.Client, ctx context.Context) {

	engine.GET("/albums", getAlbums)
	engine.GET("/albums/:id", getAlbumByID)
	engine.POST("/album", postAlbum)

	// firebase endpoints
	engine.GET("/todos", getTodos(client, ctx))
	engine.POST("/todos", createTodos(client, ctx))
}
