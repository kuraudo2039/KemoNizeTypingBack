package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	apiTest "gin_test/api/test"
	util "gin_test/util"

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
	}))

	// endpoints
	engine.GET("/helloworld", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hello world",
		})
	})

	// regist test api
	apiTest.RegistApi(engine, client, ctx)

	engine.Run(":3000")
}
