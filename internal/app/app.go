package app

import (
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"go-pocket-link/internal/config"
	"log"
)

func Run(configPath string) {
	cfg := mustReadConfig(config.NewFileReader(configPath))
	log.Println("read file", configPath)

	postgresDB := mustConnectToPostgres(cfg)
	defer func() { _ = postgresDB.Close() }()
}

func mustReadConfig(reader config.Reader) *config.Config {
	cfg, err := reader.Read()
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}

func mustConnectToPostgres(cfg *config.Config) *sqlx.DB {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Storage.Postgres.User, cfg.Storage.Postgres.Password,
		cfg.Storage.Postgres.Host, cfg.Storage.Postgres.Port,
		cfg.Storage.Postgres.Name, cfg.Storage.Postgres.SslMode)
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connected to database", fmt.Sprintf("postgres://%s:%d/%s",
		cfg.Storage.Postgres.Host, cfg.Storage.Postgres.Port, cfg.Storage.Postgres.Name))
	return db
}
