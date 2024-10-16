package service

import (
	"context"
	"go-pocket-link/pkg/email"
	"go-pocket-link/pkg/errb"
)

type EmailService struct {
	dialer email.Dialer
}

func NewEmailService(dialer email.Dialer) *EmailService {
	return &EmailService{dialer: dialer}
}

func (s *EmailService) Send(ctx context.Context, messages ...email.Message) error {
	if len(messages) > 1 {
		return s.dialer.SendAll(ctx, messages)
	} else if len(messages) > 0 {
		return s.dialer.Send(ctx, messages[0])
	}
	return errb.Errorf("no messages provided")
}
