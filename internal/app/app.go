package app

import (
	"database/sql"
	_ "github.com/lib/pq"
	"go-pocket-link/internal/config"
	"log"
)

func Run(configPath string) {
	cfg := mustReadConfig(config.NewFileReader(configPath))

	dsn := cfg.DB.Postgres.DSN()
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = db.Close() }()

	log.Println("connected to", dsn)
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connection is stable")
}

func mustReadConfig(r config.Reader) *config.Config {
	cfg := r.MustRead()
	return cfg
}
