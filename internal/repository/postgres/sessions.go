package postgres

import (
	"context"
	"github.com/google/uuid"
	"go-pocket-link/internal/domain"
	"go-pocket-link/pkg/database/postgres"
)

type SessionsRepository struct {
	db *postgres.DB
}

func NewSessionsRepository(db *postgres.DB) *SessionsRepository {
	return &SessionsRepository{db: db}
}

func (r *SessionsRepository) Save(ctx context.Context, session *domain.Session) error {
	return r.db.Save(ctx, &session.ID, `INSERT INTO sessions (user_id, refresh_token) VALUES (:user_id, :refresh_token) RETURNING id`, session)
}

func (r *SessionsRepository) Get(ctx context.Context, id uuid.UUID) (domain.Session, error) {
	var session domain.Session
	if err := r.db.GetPrepared(ctx, `SELECT * FROM sessions WHERE id = $1`, id.String()); err != nil {
		return domain.Session{}, err
	}
	return session, nil
}

func (r *SessionsRepository) GetByUserID(ctx context.Context, id uuid.UUID) (domain.Session, error) {
	var session domain.Session
	if err := r.db.GetPrepared(ctx, `SELECT * FROM sessions WHERE user_id = $1`, id.String()); err != nil {
		return domain.Session{}, err
	}
	return session, nil
}

func (r *SessionsRepository) GetByRefreshToken(ctx context.Context, token string) (domain.Session, error) {
	var session domain.Session
	if err := r.db.GetPrepared(ctx, `SELECT * FROM sessions WHERE refresh_token = $1`, token); err != nil {
		return domain.Session{}, err
	}
	return session, nil
}

func (r *SessionsRepository) Update(ctx context.Context, session *domain.Session) error {
	return r.db.UpdateNamed(ctx, `UPDATE sessions SET refresh_token = :refresh_token, is_invoked = :is_invoked WHERE id = :id`, session)
}

func (r *SessionsRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.Delete(ctx, `DELETE FROM sessions WHERE id = $1`, id.String())
}
