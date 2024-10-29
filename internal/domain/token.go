package domain

import (
	"fmt"
	"github.com/google/uuid"
)

type Token struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
}

func (t *Token) Key() string {
	return fmt.Sprintf("%s:%s", t.UserID.String(), t.ID.String())
}
