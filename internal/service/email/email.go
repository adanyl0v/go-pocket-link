package email

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/go-mail/mail"
	"go-pocket-link/internal/service"
)

type notifierImpl struct {
	dialer *mail.Dialer
}

type DialerOptions struct {
	Username  string
	Password  string
	TLSConfig *tls.Config
}

const (
	TypeText = "text/plain"
	TypeHTML = "text/html"
)

type Message struct {
	From    string   `json:"from"`
	To      string   `json:"to"`
	Cc      []string `json:"cc,omitempty"`
	Bcc     []string `json:"bcc,omitempty"`
	Subject string   `json:"subject"`
	Type    string   `json:"type"`
	Body    string   `json:"body"`
}

func NewNotifier(dialerOpts *DialerOptions) (service.EmailNotifier, error) {
	if dialerOpts == nil {
		return nil, fmt.Errorf("no dialer options provided")
	}
	dialer := mail.NewDialer("smtp.gmail.com", 587, dialerOpts.Username, dialerOpts.Password)
	if dialerOpts.TLSConfig != nil {
		dialer.TLSConfig = dialerOpts.TLSConfig
	}
	return &notifierImpl{dialer: dialer}, nil
}

func (n *notifierImpl) DialerUsername() string {
	return n.dialer.Username
}

func (n *notifierImpl) Send(ctx context.Context, message Message) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		err := n.dialer.DialAndSend(n.newMassage(&message))
		if err != nil {
			return err
		}
		return nil
	}
}

func (n *notifierImpl) newMassage(message *Message) *mail.Message {
	m := mail.NewMessage()
	m.SetHeader("From", message.From)
	m.SetHeader("To", message.To)
	if message.Cc != nil && len(message.Cc) > 0 {
		m.SetHeader("Cc", message.Cc...)
	}
	if message.Bcc != nil && len(message.Bcc) > 0 {
		m.SetHeader("Bcc", message.Bcc...)
	}
	m.SetHeader("Subject", message.Subject)
	m.SetBody(message.Type, message.Body)
	return m
}
