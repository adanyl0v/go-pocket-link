package app

import (
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go-pocket-link/internal/config"
	"go-pocket-link/pkg/database/postgres"
	"log"
	"log/slog"
	"os"
)

func Run(configPath string) {
	cfg := mustReadConfig(config.NewFileReader(configPath))
	log.Println("read config file", configPath)

	mustSetupLogger(cfg.Env)
	slog.Info("set up logger", "env", cfg.Env)

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

func mustSetupLogger(env string) {
	var logger *slog.Logger
	switch env {
	case config.EnvLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		}))
	case config.EnvProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelWarn,
		}))
	case config.EnvDev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		}))
	default:
		log.Fatalln("unknown env", env)
	}
	slog.SetDefault(logger)
}

func mustConnectToPostgres(cfg *config.Config) *postgres.DB {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Storage.Postgres.User, cfg.Storage.Postgres.Password,
		cfg.Storage.Postgres.Host, cfg.Storage.Postgres.Port,
		cfg.Storage.Postgres.Name, cfg.Storage.Postgres.SslMode)
	db, err := postgres.Connect(dsn, &postgres.ConnOptions{
		MaxOpenConns:    cfg.Storage.Postgres.MaxOpenConns,
		MaxIdleConns:    cfg.Storage.Postgres.MaxIdleConns,
		ConnMaxLifetime: cfg.Storage.Postgres.ConnMaxLifetime,
		ConnMaxIdleTime: cfg.Storage.Postgres.ConnMaxIdleTime,
	})
	if err != nil {
		slog.Error("connecting to postgres", "err", err.Error())
		os.Exit(1)
	}
	slog.Info("connected to database", "dsn", fmt.Sprintf("postgres://%s:%d/%s",
		cfg.Storage.Postgres.Host, cfg.Storage.Postgres.Port, cfg.Storage.Postgres.Name))
	return db
}
