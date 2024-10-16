package service

import (
	"context"
	"go-pocket-link/internal/service/email"
)

type EmailNotifier interface {
	DialerUsername() string
	Send(ctx context.Context, message email.Message) error
}
