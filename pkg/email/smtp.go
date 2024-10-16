package email

import (
	"context"
	"crypto/tls"
	"github.com/go-mail/mail"
)

type SMTP struct {
	dialer *mail.Dialer
}

func NewSMTPDialer(username, password string, config *tls.Config) *SMTP {
	d := mail.NewDialer("smtp.gmail.com", 587, username, password)
	d.TLSConfig = config
	return &SMTP{d}
}

func (s *SMTP) Send(ctx context.Context, message Message) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		err := s.dialer.DialAndSend(s.createMessage(&message))
		if err != nil {
			return err
		}
		return nil
	}
}

func (s *SMTP) SendAll(ctx context.Context, messages []Message) error {
	mailMessages := make([]*mail.Message, 0, cap(messages))
	for _, message := range messages {
		mailMessages = append(mailMessages, s.createMessage(&message))
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		err := s.dialer.DialAndSend(mailMessages...)
		if err != nil {
			return err
		}
		return nil
	}
}

func (s *SMTP) createMessage(message *Message) *mail.Message {
	m := mail.NewMessage()
	if message.From != "" {
		m.SetHeader("From", message.From)
	} else {
		m.SetHeader("From", s.dialer.Username)
	}
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
