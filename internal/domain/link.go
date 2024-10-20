package domain

import (
	"github.com/google/uuid"
	"time"
)

type Link struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Title     string    `json:"title" db:"title"`
	URL       string    `json:"url" db:"url"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
