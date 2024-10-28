package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go-pocket-link/internal/config"
	delivhttp "go-pocket-link/internal/delivery/http"
	httpv1 "go-pocket-link/internal/delivery/http/v1"
	"go-pocket-link/internal/repository"
	pgrep "go-pocket-link/internal/repository/postgres"
	pgdb "go-pocket-link/pkg/database/postgres"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	logError = "error"
)

func Run(configPath string) {
	cfg := mustReadConfig(config.NewFileReader(configPath))
	log.Println("read config file", configPath)

	mustSetupLogger(cfg.Env)
	slog.Info("set up logger", "env", cfg.Env)

	postgresDB := mustConnectToPostgres(cfg)
	defer func() { _ = postgresDB.Close() }()

	repos := &repository.Repositories{
		Users: pgrep.NewUsersRepository(postgresDB),
	}
	_ = repos //TODO delete me

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	delivhttp.InitRouter(router, httpv1.NewHandler())

	server := http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}
	mustListenAndServe(&server)
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

func mustConnectToPostgres(cfg *config.Config) *pgdb.DB {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Storage.Postgres.User, cfg.Storage.Postgres.Password,
		cfg.Storage.Postgres.Host, cfg.Storage.Postgres.Port,
		cfg.Storage.Postgres.Name, cfg.Storage.Postgres.SslMode)
	db, err := pgdb.Connect(dsn, &pgdb.ConnOptions{
		MaxOpenConns:    cfg.Storage.Postgres.MaxOpenConns,
		MaxIdleConns:    cfg.Storage.Postgres.MaxIdleConns,
		ConnMaxLifetime: cfg.Storage.Postgres.ConnMaxLifetime,
		ConnMaxIdleTime: cfg.Storage.Postgres.ConnMaxIdleTime,
	})
	if err != nil {
		slog.Error("connecting to postgres", logError, err.Error())
		os.Exit(1)
	}
	slog.Info("connected to database", "dsn", fmt.Sprintf("postgres://%s:%d/%s",
		cfg.Storage.Postgres.Host, cfg.Storage.Postgres.Port, cfg.Storage.Postgres.Name))
	return db
}

func mustListenAndServe(server *http.Server) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("", logError, err)
			os.Exit(1)
		}
		slog.Info("listening...", "addr", server.Addr)
	}()

	<-ctx.Done()
	slog.Info("shutting down server...")
	if err := server.Shutdown(context.Background()); err != nil {
		slog.Error("", logError, err)
		os.Exit(1)
	}
	slog.Info("shut down gracefully")
}
