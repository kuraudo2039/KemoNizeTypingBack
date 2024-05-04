package errorObj

import "github.com/gin-gonic/gin"

func CreateErr(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func CreateErrFromString(msg string) gin.H {
	return gin.H{"error": msg}
}
