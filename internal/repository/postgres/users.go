package postgres

import (
	"context"
	"github.com/google/uuid"
	"go-pocket-link/internal/domain"
	"go-pocket-link/pkg/database/postgres"
)

type UsersRepository struct {
	db *postgres.DB
}

func NewUsersRepository(db *postgres.DB) *UsersRepository {
	return &UsersRepository{db: db}
}

func (r *UsersRepository) Save(ctx context.Context, user *domain.User) error {
	err := r.db.Save(ctx, &user.ID, `INSERT INTO users(name, email, password) VALUES (:name, :email, :password) RETURNING id`, user)
	if err != nil {
		return err
	}
	return nil
}

func (r *UsersRepository) Get(ctx context.Context, id uuid.UUID) (domain.User, error) {
	var user domain.User
	err := r.db.GetPrepared(ctx, &user, `SELECT * FROM users WHERE id = $1`, id.String())
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (r *UsersRepository) GetByCredentials(ctx context.Context, email, password string) (domain.User, error) {
	var user domain.User
	err := r.db.GetPrepared(ctx, &user, `SELECT * FROM users WHERE email = $1 AND password = $2`, email, password)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (r *UsersRepository) Update(ctx context.Context, user *domain.User) error {
	err := r.db.Update(ctx, `UPDATE users SET name = :name, email = :email, password = :password WHERE id = :id`, user)
	if err != nil {
		return err
	}
	return nil
}

func (r *UsersRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.db.Delete(ctx, `DELETE FROM users WHERE id = $1`, id.String())
	if err != nil {
		return err
	}
	return nil
}
