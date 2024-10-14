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
	query := `INSERT INTO users (name, email, password, created_at, updated_at) VALUES (:name, :email, :password, :created_at, :updated_at) RETURNING id`
	user.CreatedAt = time.Now().In(time.Local)
	user.UpdatedAt = user.CreatedAt
	return r.db.Save(ctx, &user.ID, query, user)
}

func (r *UsersRepository) GetByID(ctx context.Context, dest *domain.User) error {
	query := `SELECT * FROM users WHERE id = :id`
	return r.db.GetNamed(ctx, dest, query, dest)
}

func (r *UsersRepository) GetByCredentials(ctx context.Context, dest *domain.User) error {
	query := `SELECT * FROM users WHERE email = :email AND password = :password`
	return r.db.GetNamed(ctx, dest, query, dest)
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
	query := `UPDATE users SET name = :name, email = :email, password = :password, updated_at = :updated_at WHERE id = :id`
	user.UpdatedAt = time.Now().In(time.Local)
	return r.db.Update(ctx, query, user)
}

func (r *UsersRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	return r.db.Delete(ctx, query, id)
}
