package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-pocket-link/pkg/email"
	"go-pocket-link/pkg/errb"
	"net/http"
)

type emailSendInput struct {
	To  string   `json:"to"`
	Cc  []string `json:"cc,omitempty"`
	Bcc []string `json:"bcc,omitempty"`
}

//TODO remove raw body strings and replace them with HTML templates

func (h *handlerImpl) EmailSendSuccessfulSignIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.sendEmail(c, "Welcome back to Pocket Link!", `
		<html>
			<body>
				<header>
					<h1>You've signed in successfully</h1>
				</header>
			</body>
		</html>
		`)
	}
}

func (h *handlerImpl) EmailSendSuccessfulSignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.sendEmail(c, "Welcome to Pocket Link!", `
		<html>
			<body>
				<header>
					<p>You've signed up successfully</p>
				</header>
			</body>
		</html>
		`)
	}
}

func (h *handlerImpl) sendEmail(c *gin.Context, subject, body string) {
	message := emailSendInput{}
	err := c.ShouldBindJSON(&message)
	if err != nil {
		writeError(c, http.StatusBadRequest, errb.Errorf("binding json: %v", err))
		return
	}

	err = h.services.Email.Send(c, email.Message{
		To:      message.To,
		Cc:      message.Cc,
		Bcc:     message.Bcc,
		Subject: subject,
		Type:    email.TypeTextHTML,
		Body:    body,
	})
	if err != nil {
		writeError(c, http.StatusInternalServerError, errb.Errorf("sending an email: %v", err))
		return
	}

	writeMessage(c, http.StatusOK, gin.H{keyMessage: fmt.Sprintln("email sent to", message.To)})
}
