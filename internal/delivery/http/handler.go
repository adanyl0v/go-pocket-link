package http

import (
	"github.com/gin-gonic/gin"
)

type EndpointsInitializer interface {
	InitEndpoints(group *gin.RouterGroup)
}

func InitRouter(router *gin.Engine, v1Handler EndpointsInitializer) {
	v1Group := router.Group("/api/v1")
	v1Handler.InitEndpoints(v1Group)
}
