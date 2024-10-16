package email

import (
	"context"
)

type Dialer interface {
	Send(ctx context.Context, message Message) error
	SendAll(ctx context.Context, messages []Message) error
}

type Message struct {
	From    string   `json:"from"`
	To      string   `json:"to"`
	Cc      []string `json:"cc,omitempty"`
	Bcc     []string `json:"bcc,omitempty"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
	Type    string   `json:"type"`
}

const (
	TypeTextPlain = "text/plain"
	TypeTextHTML  = "text/html"
)
