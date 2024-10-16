package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go-pocket-link/internal/service"
	"go-pocket-link/pkg/email"
	"net/http"
)

type EmailHandler struct {
	service *service.EmailService
}

func NewEmailHandler(email *service.EmailService) *EmailHandler {
	return &EmailHandler{service: email}
}

func (h *EmailHandler) Send() gin.HandlerFunc {
	return func(c *gin.Context) {
		inp := email.Message{}
		err := c.ShouldBindJSON(&inp)
		if err != nil {
			err = fmt.Errorf("failed to parse a json email input: %v", err)
			log.Errorln(err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		switch inp.Type {
		case email.TypeTextPlain:
		case email.TypeTextHTML:
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid email content type",
			})
		}

		err = h.service.Send(c, email.Message{
			To:      inp.To,
			Cc:      inp.Cc,
			Bcc:     inp.Bcc,
			Subject: inp.Subject,
			Type:    inp.Type,
			Body:    inp.Body,
		})
		if err != nil {
			err = fmt.Errorf("failed to send an email: %v", err)
			log.Errorln(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
		}
	}
}
