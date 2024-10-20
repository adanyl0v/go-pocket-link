package v1

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-pocket-link/pkg/email"
	"go-pocket-link/pkg/errb"
	"html/template"
	"net/http"
	"path"
)

type emailSendInput struct {
	To  string   `json:"to"`
	Cc  []string `json:"cc,omitempty"`
	Bcc []string `json:"bcc,omitempty"`
}

func (h *handlerImpl) EmailSendSuccessfulSignIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		filePath := path.Join(h.services.Email.TemplatesDir, "successful_sign_in.html")
		b, err := htmlTemplateToBytesBuffer(filePath, nil)
		if err != nil {
			writeError(c, http.StatusInternalServerError, err)
			return
		}
		h.sendEmail(c, "Welcome back to Pocket Link!", b.String())
	}
}

func (h *handlerImpl) EmailSendSuccessfulSignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		filePath := path.Join(h.services.Email.TemplatesDir, "successful_sign_up.html")
		b, err := htmlTemplateToBytesBuffer(filePath, nil)
		if err != nil {
			writeError(c, http.StatusInternalServerError, err)
			return
		}
		h.sendEmail(c, "Welcome to Pocket Link!", b.String())
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

	writeMessage(c, http.StatusOK, gin.H{keyMessage: fmt.Sprint("email sent to", message.To)})
}

func htmlTemplateToBytesBuffer(path string, data any) (*bytes.Buffer, error) {
	t, err := template.ParseFiles(path)
	if err != nil {
		return nil, errb.Errorf("parsing HTML template: %v", err)
	}
	b := bytes.NewBufferString("")
	err = t.Execute(b, data)
	if err != nil {
		return nil, errb.Errorf("executing HTML template: %v", err)
	}
	return b, nil
}
