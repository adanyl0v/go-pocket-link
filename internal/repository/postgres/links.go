package postgres

import (
	"context"
	"github.com/google/uuid"
	"go-pocket-link/internal/domain"
	"go-pocket-link/pkg/storage"
	"time"
)

type LinksRepository struct {
	db storage.DB
}

func NewLinksRepository(db storage.DB) *LinksRepository {
	return &LinksRepository{db: db}
}

func (r *LinksRepository) Save(ctx context.Context, link *domain.Link) error {
	query := `INSERT INTO links (title, url, user_id) VALUES (:title, :url, :user_id)`
	return r.db.Save(ctx, &link.ID, query, link)
}

func (r *LinksRepository) GetByID(ctx context.Context, dest *domain.Link) error {
	query := `SELECT * FROM links WHERE id = $1`
	return r.db.Get(ctx, dest, query, dest.ID)
}

func (r *LinksRepository) GetByTitle(ctx context.Context, title string) ([]domain.Link, error) {
	query := `SELECT * FROM links WHERE title = $1`
	var links []domain.Link
	err := r.db.GetAll(ctx, &links, query, title)
	if err != nil {
		return nil, err
	}
	return links, nil
}

func (r *LinksRepository) GetByURL(ctx context.Context, dest *domain.Link) error {
	query := `SELECT * FROM links WHERE url = $1`
	return r.db.Get(ctx, dest, query, dest.ID)
}

func (r *LinksRepository) GetAll(ctx context.Context) ([]domain.Link, error) {
	query := `SELECT * FROM links`
	var links []domain.Link
	err := r.db.GetAll(ctx, &links, query)
	if err != nil {
		return nil, err
	}
	return links, nil
}

func (r *LinksRepository) GetAllByUserID(ctx context.Context, id uuid.UUID) ([]domain.Link, error) {
	query := `SELECT * FROM links WHERE user_id = $1`
	var links []domain.Link
	err := r.db.GetAll(ctx, &links, query, id)
	if err != nil {
		return nil, err
	}
	return links, nil
}

func (r *LinksRepository) Update(ctx context.Context, link *domain.Link) error {
	query := `UPDATE links SET title = :title, url = :url, updated_at = :updated_at WHERE id = :id`
	link.UpdatedAt = time.Now().UTC()
	return r.db.Save(ctx, &link.ID, query, link)
}

func (r *LinksRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM links WHERE id = $1`
	return r.db.Delete(ctx, query, id)
}
