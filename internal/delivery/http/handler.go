package http

import (
	"github.com/gin-gonic/gin"
	"go-pocket-link/internal/delivery/http/v1"
	"go-pocket-link/internal/service"
)

type EndpointsInitializer interface {
	InitEndpoints(parent *gin.RouterGroup)
}

type Handler struct {
	v1Handler EndpointsInitializer
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{v1Handler: v1.NewHandler(services)}
}

func (h *Handler) InitRoutes(router *gin.Engine) {
	apiGroupV1 := router.Group("/api/v1")
	h.v1Handler.InitEndpoints(apiGroupV1)
}
