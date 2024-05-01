package errorObj

import "github.com/gin-gonic/gin"

func CreateErr(err error) gin.H {
	return gin.H{"error": err.Error()}
}
