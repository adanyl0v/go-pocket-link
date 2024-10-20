package postgres

import (
	"context"
	"github.com/google/uuid"
	"go-pocket-link/internal/domain"
	"go-pocket-link/pkg/storage"
	"time"
)

type UsersRepository struct {
	db storage.DB
}

func NewUsersRepository(db storage.DB) *UsersRepository {
	return &UsersRepository{db: db}
}

func (r *UsersRepository) Save(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (name, email, password, created_at) VALUES (:name, :email, :password, :created_at) RETURNING id`
	user.CreatedAt = time.Now().UTC()
	return r.db.Save(ctx, &user.ID, query, user)
}

func (r *UsersRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	var user domain.User
	query := `SELECT * FROM users WHERE id = $1`
	err := r.db.Get(ctx, &user, query, id)
	return user, err
}

func (r *UsersRepository) GetByCredentials(ctx context.Context, email, password string) (domain.User, error) {
	var user domain.User
	query := `SELECT * FROM users WHERE email = $1 AND password = $2`
	err := r.db.Get(ctx, &user, query, email, password)
	return user, err
}

func (r *UsersRepository) GetAll(ctx context.Context) ([]domain.User, error) {
	query := `SELECT * FROM users`
	var users []domain.User
	err := r.db.GetAll(ctx, &users, query)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UsersRepository) Update(ctx context.Context, user *domain.User) error {
	query := `UPDATE users SET name = :name, email = :email, password = :password WHERE id = :id`
	return r.db.Update(ctx, query, user)
}

func (r *UsersRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	return r.db.Delete(ctx, query, id)
}
