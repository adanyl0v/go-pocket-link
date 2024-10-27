package database

import "context"

type DB interface {
	Close() error

	Saver
	Getter
	Selector
	Updater
	Deleter
}

type Saver interface {
	Save(ctx context.Context, dest any, query string, arg any) error
}

type Getter interface {
	Get(ctx context.Context, dest any, query string, args ...any) error
	GetNamed(ctx context.Context, dest any, query string, arg any) error
	GetPrepared(ctx context.Context, dest any, query string, args ...any) error
}

type Selector interface {
	Select(ctx context.Context, dest any, query string, args ...any) error
	SelectPrepared(ctx context.Context, dest any, query string, args ...any) error
}

type Updater interface {
	Update(ctx context.Context, query string, args ...any) error
	UpdateNamed(ctx context.Context, query string, arg any) error
}

type Deleter interface {
	Delete(ctx context.Context, query string, args ...any) error
	DeleteNamed(ctx context.Context, query string, arg any) error
}
