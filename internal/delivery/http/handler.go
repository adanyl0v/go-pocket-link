package http

import (
	"github.com/gin-gonic/gin"
	"go-pocket-link/internal/delivery/http/v1"
	"go-pocket-link/internal/service"
	"net/http"
)

type Handler struct {
	emailService *service.EmailService
}

func NewHandler(email *service.EmailService) *Handler {
	return &Handler{emailService: email}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	api.GET("/ping", h.handlePing())

	emailGroup := api.Group("/email")
	{
		handler := v1.NewEmailHandler(h.emailService)
		emailGroup.POST("/send", handler.Send())
	}
}

func (h *Handler) handlePing() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	}
}
