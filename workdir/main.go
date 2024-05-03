package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	apiTest "gin_test/coyote/api/test"
	coyoteWsApi "gin_test/coyote/api/ws"
	coyoteHttpApi "gin_test/coyote/http"
	util "gin_test/coyote/util"

	"github.com/gin-contrib/cors"
)

// main
func main() {

	// initFirebaseClient
	client, ctx := util.InitFirebaseClient()
	defer client.Close()

	// initEngine
	engine := gin.Default()
	engine.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Type"},
	}))

	// endpoints
	engine.GET("/helloworld", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hello world",
		})
	})

	// regist test api
	apiTest.RegistApi(engine, client, ctx)
	go apiTest.WsHandleMessages()

	// regist coyote api
	coyoteHttpApi.RegistHttpApi(engine, client)
	go coyoteWsApi.HandleMessages()

	engine.Run(":3000")
}
