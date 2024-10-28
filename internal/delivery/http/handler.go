package http

import (
	"github.com/gin-gonic/gin"
	v1 "go-pocket-link/internal/delivery/http/v1"
)

func InitRouter(router *gin.Engine, v1Handler v1.Handler) {
	v1Group := router.Group("/api/v1")
	v1Handler.InitEndpoints(v1Group)
}
