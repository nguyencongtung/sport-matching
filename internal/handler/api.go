package handler

import "github.com/gin-gonic/gin"

// Hello handle api status
func Hello(c *gin.Context) {
	c.JSON(200, gin.H{"status": "success", "message": "Hello i'm ok!", "data": nil})
}
