package v1

import (
	"github.com/gin-gonic/gin"
	"log/slog"
)

func writeError(c *gin.Context, status int, message string, err error) {
	c.JSON(status, gin.H{logError: message})
	if err == nil {
		slog.Error(message)
	} else {
		slog.Error(message, logError, err)
	}
}

func writeAbort(c *gin.Context, status int, message string, err error) {
	writeError(c, status, message, err)
	c.Abort()
}
