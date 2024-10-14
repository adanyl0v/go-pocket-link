package http

import (
	"github.com/gin-gonic/gin"
	"go-pocket-link/internal/delivery/http/v1"
	"go-pocket-link/internal/service"
)

type Handler struct {
	emailNotifier service.EmailNotifier
}

func NewHandler(emailNotifier service.EmailNotifier) *Handler {
	return &Handler{
		emailNotifier: emailNotifier,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	emailGroup := api.Group("/email")
	{
		handler := v1.NewEmailHandler(h.emailNotifier)
		emailGroup.POST("/send", handler.Send())
	}
}
