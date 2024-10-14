package service

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/go-mail/mail"
)

type EmailNotifier interface {
	DialerUsername() string
	AddMessage(message EmailMessage)
	Send(ctx context.Context) error
}

type emailNotifierImpl struct {
	dialer   *mail.Dialer
	messages []*mail.Message
}

type EmailDialerOptions struct {
	Username  string
	Password  string
	TLSConfig *tls.Config
}

const (
	EmailTypeText = "text/plain"
	EmailTypeHTML = "text/html"
)

type EmailMessage struct {
	From    string   `json:"from"`
	To      string   `json:"to"`
	Cc      []string `json:"cc,omitempty"`
	Bcc     []string `json:"bcc,omitempty"`
	Subject string   `json:"subject"`
	Type    string   `json:"type"`
	Body    string   `json:"body"`
}

func NewEmailNotifier(dialerOpts *EmailDialerOptions) (EmailNotifier, error) {
	if dialerOpts == nil {
		return nil, fmt.Errorf("no dialer options provided")
	}
	dialer := mail.NewDialer("smtp.gmail.com", 587, dialerOpts.Username, dialerOpts.Password)
	if dialerOpts.TLSConfig != nil {
		dialer.TLSConfig = dialerOpts.TLSConfig
	}
	return &emailNotifierImpl{dialer: dialer}, nil
}

func (n *emailNotifierImpl) DialerUsername() string {
	return n.dialer.Username
}

func (n *emailNotifierImpl) AddMessage(message EmailMessage) {
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
	n.messages = append(n.messages, m)
}

func (n *emailNotifierImpl) Send(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		err := n.dialer.DialAndSend(n.messages...)
		if err != nil {
			return err
		}
		n.messages = n.messages[:0]
		return nil
	}
}
