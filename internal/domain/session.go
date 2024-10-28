package domain

import (
	"github.com/google/uuid"
	"time"
)

type Session struct {
	ID           uuid.UUID `json:"id" db:"id"`
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	RefreshToken string    `json:"refresh_token" db:"refresh_token"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	IsInvoked    bool      `json:"is_invoked" db:"is_invoked"`
}
