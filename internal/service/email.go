package service

import (
	"context"
	"go-pocket-link/pkg/email"
	"go-pocket-link/pkg/errb"
)

type EmailService struct {
	dialer       email.Dialer
	TemplatesDir string
}

func NewEmailService(dialer email.Dialer, templatesDir string) *EmailService {
	return &EmailService{
		dialer:       dialer,
		TemplatesDir: templatesDir,
	}
}

func (s *EmailService) Send(ctx context.Context, messages ...email.Message) error {
	if len(messages) > 1 {
		return s.dialer.SendAll(ctx, messages)
	} else if len(messages) > 0 {
		return s.dialer.Send(ctx, messages[0])
	}
	return errb.Errorf("no messages provided")
}
