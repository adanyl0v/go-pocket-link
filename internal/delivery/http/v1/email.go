package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go-pocket-link/internal/service"
	"go-pocket-link/internal/service/email"
	"net/http"
)

type EmailHandler struct {
	notifier service.EmailNotifier
}

func NewEmailHandler(notifier service.EmailNotifier) *EmailHandler {
	return &EmailHandler{
		notifier: notifier,
	}
}

type emailInput struct {
	To      string   `json:"to"`
	Cc      []string `json:"cc,omitempty"`
	Bcc     []string `json:"bcc,omitempty"`
	Subject string   `json:"subject"`
	Type    string   `json:"type"`
	Body    string   `json:"body"`
}

func (h *EmailHandler) Send() gin.HandlerFunc {
	return func(c *gin.Context) {
		inp := emailInput{}
		err := c.ShouldBindJSON(&inp)
		if err != nil {
			err = fmt.Errorf("failed to parse a json email input: %v", err)
			log.Errorln(err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
		err = h.notifier.Send(c, email.Message{
			From:    h.notifier.DialerUsername(),
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
