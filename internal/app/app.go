package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/adanyl0v/go-pocket-link/internal/config"
	delivhttp "github.com/adanyl0v/go-pocket-link/internal/delivery/http"
	httpv1 "github.com/adanyl0v/go-pocket-link/internal/delivery/http/v1"
	"github.com/adanyl0v/go-pocket-link/internal/repository"
	pgrep "github.com/adanyl0v/go-pocket-link/internal/repository/postgres"
	redisrep "github.com/adanyl0v/go-pocket-link/internal/repository/redis"
	"github.com/adanyl0v/go-pocket-link/internal/service"
	"github.com/adanyl0v/go-pocket-link/pkg/auth/jwt"
	redisdb "github.com/adanyl0v/go-pocket-link/pkg/cache/redis"
	"github.com/adanyl0v/go-pocket-link/pkg/crypto/hash"
	pgdb "github.com/adanyl0v/go-pocket-link/pkg/database/postgres"
	"github.com/adanyl0v/go-pocket-link/pkg/validator"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
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

	mustSetupLogger(cfg.Env)
	slog.Info("read config", "path", configPath)
	slog.Info("set up logger", "env", cfg.Env)

	postgresDB := mustConnectToPostgres(cfg)
	defer func() { _ = postgresDB.Close() }()

	redisDB := mustConnectToRedis(cfg)
	defer func() { _ = redisDB.Close() }()

	repos := &repository.Repositories{
		Users:  pgrep.NewUsersRepository(postgresDB),
		Tokens: redisrep.NewTokensRepository(redisDB),
	}
	slog.Info("initialized repositories")

	services := service.Services{
		Users: service.NewUsersService(repos.Users, hash.NewSHA1Hasher(cfg.Hash.Salt), validator.NewCredentialsValidator()),
		Tokens: service.NewTokensService(repos.Tokens, jwt.NewTokenManager(cfg.Auth.AccessSecret, cfg.Auth.RefreshSecret,
			jwt.StaticClaims{Issuer: "https://pocketlink.com", Audience: "https://api.pocketlink.com"}),
			cfg.Auth.AccessTokenTTL, cfg.Auth.RefreshTokenTTL),
	}
	slog.Info("initialized services")

	router := gin.New()
	router.Use(gin.Recovery())

	mustSetupRouterLogger(router, cfg.Env)
	//TODO: how about adding ELK support?

	delivhttp.InitRouter(router, httpv1.NewHandler(&services))
	slog.Info("initialized router")

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

func mustSetupRouterLogger(router *gin.Engine, env string) {
	var logger *slog.Logger

	switch env {
	case config.EnvLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug,
			AddSource: true})).With("env", env).With("mode", os.Getenv(gin.EnvGinMode))
	case config.EnvProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo,
			AddSource: true})).With("env", env).With("mode", os.Getenv(gin.EnvGinMode))
	case config.EnvDev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug,
			AddSource: true})).With("env", env).With("mode", os.Getenv(gin.EnvGinMode))
	}

	loggerConfig := sloggin.Config{
		DefaultLevel:     slog.LevelDebug,
		ClientErrorLevel: slog.LevelWarn,
		ServerErrorLevel: slog.LevelError,
		WithUserAgent:    true,
		WithRequestBody:  true,
		WithResponseBody: true,
	}

	router.Use(sloggin.NewWithConfig(logger, loggerConfig))
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
		slog.Error("connecting to postgres", logError, err)
		os.Exit(1)
	}

	slog.Info("connected to postgres", "dsn", fmt.Sprintf("postgres://...@%s:%d/%s",
		cfg.Storage.Postgres.Host, cfg.Storage.Postgres.Port, cfg.Storage.Postgres.Name))
	return db
}

func mustConnectToRedis(cfg *config.Config) *redisdb.DB {
	dsn := fmt.Sprintf("redis://:%s@%s:%d/0", cfg.Storage.Redis.Password,
		cfg.Storage.Redis.Host, cfg.Storage.Redis.Port)

	db, err := redisdb.Connect(dsn)
	if err != nil {
		slog.Error("connecting to redis", logError, err)
		os.Exit(1)
	}

	slog.Info("connected to redis", "dsn", fmt.Sprintf("redis://...@%s:%d/0",
		cfg.Storage.Redis.Host, cfg.Storage.Redis.Port))
	return db
}

func mustListenAndServe(server *http.Server) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		slog.Info("listening...", "addr", server.Addr)
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("", logError, err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	slog.Info("shutting down server...")
	if err := server.Shutdown(context.Background()); err != nil {
		slog.Error("", logError, err)
		os.Exit(1)
	}
	slog.Info("shut down gracefully")
}
