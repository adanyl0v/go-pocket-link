package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler interface {
	InitEndpoints(*gin.RouterGroup)

	Ping(*gin.Context)
}

type handlerImpl struct {
}

func NewHandler() Handler {
	return &handlerImpl{}
}

func (h *handlerImpl) InitEndpoints(routerGroup *gin.RouterGroup) {
	routerGroup.GET("/ping", h.Ping)
}

func (h *handlerImpl) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
