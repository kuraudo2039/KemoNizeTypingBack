package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// schema
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// init table
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

// handler
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	for _, item := range albums {
		if item.ID == id {
			c.IndentedJSON(http.StatusOK, item)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func postAlbum(c *gin.Context) {
	var newAlbum album

	// 受け取ったJSONを`newAlbum`にバインドするために`BindJSON`を呼び出す
	if err := c.BindJSON(&newAlbum); err != nil {
		fmt.Print(err)
	}

	// スライスへ新しいアルバムを追加する
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// firestoreクライアント初期化
func initFirebaseClient() (*firestore.Client, context.Context) {
	// Use a service account
	ctx := context.Background()
	sa := option.WithCredentialsJSON([]byte(os.Getenv("FIREBASE_CREDENTIALS")))
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	return client, ctx
}

// get todos
func getTodos(client *firestore.Client, ctx context.Context) func(c *gin.Context) {
	return func(c *gin.Context) {
		iter := client.Collection("todos").OrderBy("createdAt", firestore.Asc).Documents(ctx)

		type fireStoreCollection struct {
			ID        string                 `json:"id"`
			Documents map[string]interface{} `json:"documents"`
		}
		var data []fireStoreCollection

		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatalf("Failed to iterate: %v", err)
			}
			fmt.Println(doc.Ref.ID)
			fmt.Println(doc.Data())

			data = append(data, fireStoreCollection{doc.Ref.ID, doc.Data()})
		}

		c.IndentedJSON(http.StatusOK, data)
	}
}

// main
func main() {

	client, ctx := initFirebaseClient()
	defer client.Close()

	engine := gin.Default()

	// firebase endpoints
	engine.GET("/todos", getTodos(client, ctx))

	// endpoints
	engine.GET("/helloworld", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hello world",
		})
	})
	engine.GET("/albums", getAlbums)
	engine.GET("/albums/:id", getAlbumByID)
	engine.POST("/album", postAlbum)

	engine.Run(":3000")
}
