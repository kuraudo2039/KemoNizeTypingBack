package apiTest

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
)

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

// create todos
func createTodos(client *firestore.Client, ctx context.Context) func(c *gin.Context) {
	return func(c *gin.Context) {

		createData := map[string]interface{}{
			"content":   "request content",
			"createdAt": time.Now().Format("2006-01-02T15:04:05.000Z"),
			"done":      false,
			"title":     "request title",
		}

		resA, resB, err := client.Collection("todos").Add(ctx, createData)

		if err != nil {
			log.Fatalf("Failed to create: %v", err)
		}

		fmt.Println(resA)
		/**
		{
			"Parent": {
				"Parent": null,
				"Path": "projects/kemonizetyping/databases/(default)/documents/todos",
				"ID": "todos"
			},
			"Path": "projects/kemonizetyping/databases/(default)/documents/todos/OVS1zVTlA0IyDgjAhwAF",
			"ID": "OVS1zVTlA0IyDgjAhwAF"
		}
		*/

		fmt.Println(resB)
		/**
		{
			"UpdateTime": "2024-04-24T13:48:32.227487Z"
		}
		*/

		responseData := map[string]interface{}{
			"data":       createData,
			"response":   resA,
			"UpdateTime": resB.UpdateTime,
		}

		c.IndentedJSON(http.StatusCreated, responseData)
	}
}
