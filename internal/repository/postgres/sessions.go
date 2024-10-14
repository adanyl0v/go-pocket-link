package postgres

import (
	"context"
	"github.com/google/uuid"
	"go-pocket-link/internal/domain"
	"go-pocket-link/pkg/storage"
	"time"
)

type SessionsRepository struct {
	db storage.DB
}

func NewSessionsRepository(db storage.DB) *SessionsRepository {
	return &SessionsRepository{db: db}
}

func (r *SessionsRepository) Save(ctx context.Context, session *domain.Session) error {
	query := `INSERT INTO sessions (user_id, refresh_token, expires_at, created_at) VALUES (:user_id, :refresh_token, :expires_at, :created_at) RETURNING id`
	session.CreatedAt = time.Now().In(time.Local)
	return r.db.Save(ctx, &session.ID, query, session)
}

func (r *SessionsRepository) GetByID(ctx context.Context, dest *domain.Session) error {
	query := `SELECT * FROM sessions WHERE id = :id`
	return r.db.GetNamed(ctx, dest, query, dest)
}

func (r *SessionsRepository) GetByRefreshToken(ctx context.Context, dest *domain.Session) error {
	query := `SELECT * FROM sessions WHERE refresh_token = :refresh_token`
	return r.db.GetNamed(ctx, dest, query, dest)
}

func (r *SessionsRepository) GetAll(ctx context.Context) ([]domain.Session, error) {
	query := `SELECT * FROM sessions`
	var sessions []domain.Session
	err := r.db.GetAll(ctx, &sessions, query)
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

func (r *SessionsRepository) GetAllByUserID(ctx context.Context, id uuid.UUID) ([]domain.Session, error) {
	query := `SELECT * FROM sessions WHERE user_id = $1`
	var sessions []domain.Session
	err := r.db.GetAll(ctx, &sessions, query, id)
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

func (r *SessionsRepository) Update(ctx context.Context, session *domain.Session) error {
	query := `UPDATE sessions SET refresh_token = :refresh_token, expires_at = :expires_at WHERE id = :id`
	return r.db.Update(ctx, query, session)
}

func (r *SessionsRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM sessions WHERE id = $1`
	return r.db.Delete(ctx, query, id)
}

//func (r *SessionsRepository) DeleteByRefreshToken(ctx context.Context, token string) error {
//	query := `DELETE FROM sessions WHERE refresh_token = $1`
//	return r.db.Delete(ctx, query, token)
//}
