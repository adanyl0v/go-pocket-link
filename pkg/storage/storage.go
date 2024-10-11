package storage

import (
	"context"
)

type Result interface {
	LastInsertId() (int64, error)
}

type DB interface {
	Close() error
	DriverName() string

	Saver
	Getter
	Updater
	Deleter
}

type Saver interface {
	// Save data from arg, which must be a struct and return the last inserted id if available.
	// The result is stored in dest if dest is not nil.
	// Example: save an application user and return his new id (postgresql):
	//
	//	type User struct {
	// 		ID   int64  `db:"id"`
	// 		Name string `db:"name"`
	//	}
	// 	...
	//	user := User{ID: 0, Name: "testUser"}
	//	query := `INSERT INTO users (name) VALUES (:name) RETURNING id`
	//	err := Save(ctx, &user.ID, query, user)
	//	...
	Save(ctx context.Context, dest any, query string, arg any) error
}

type Getter interface {
	// Get matching row and store it into dest, which must be a struct
	Get(ctx context.Context, dest any, query string, args ...any) error

	// GetAll all matching rows and store them into dest, which must be a slice pointer
	GetAll(ctx context.Context, dest any, query string, args ...any) error
}

type Updater interface {
	Update(ctx context.Context, query string, arg any) error
}

type Deleter interface {
	Delete(ctx context.Context, query string, args ...any) error
}
