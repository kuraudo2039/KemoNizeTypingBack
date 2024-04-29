package apiTest

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Schema
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
