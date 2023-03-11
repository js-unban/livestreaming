package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func sendSuccessResponse(c *gin.Context, result any) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"result": result,
	})
}

func sendErrorResponse(c *gin.Context, code int, message string) {
	c.IndentedJSON(code, gin.H{
		"error":     message,
		"timestamp": time.Now().String(),
		"path":      c.Request.RequestURI,
	})
}
