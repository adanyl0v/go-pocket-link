package app

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go-pocket-link/internal/config"
	httpdeliv "go-pocket-link/internal/delivery/http"
	"go-pocket-link/internal/repository"
	"go-pocket-link/internal/service"
	"go-pocket-link/pkg/email"
	"go-pocket-link/pkg/errb"
	pgstor "go-pocket-link/pkg/storage/postgres"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
)

func Run(configPath string) {
	cfg := mustReadConfig(config.NewFileReader(configPath))
	mustSetupLogger(cfg.Env)

	log.Infof("read config '%s'\n", configPath)
	log.Infof("set up logger with '%s' env\n", cfg.Env)

	pgDB := mustSetupDB(cfg.DB)
	defer func() { _ = pgDB.Close() }()
	log.Infoln("connected to postgres")

	repos := repository.NewRepositories(pgDB)
	var httpHandler *httpdeliv.Handler
	{
		_ = repos //TODO remove me
		httpHandler = httpdeliv.NewHandler(
			service.NewEmailService(email.NewSMTPDialer(
				cfg.Email.Username,
				cfg.Email.Password,
				&tls.Config{InsecureSkipVerify: true},
			)),
		)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	httpHandler.Init(router.Group("/api/v1"))

	log.Infof("listening to %s:%d...\n", cfg.Server.Host, cfg.Server.Port)
	server := http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	mustRunServer(&server)
}

func mustReadConfig(r config.Reader) *config.Config {
	cfg := r.MustRead()
	return cfg
}

func mustSetupLogger(env string) {
	l := log.StandardLogger()

	switch env {
	case config.EnvLocal:
		l.SetLevel(log.DebugLevel)
		l.SetFormatter(&log.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "02.01.2006 15:04:05",
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				dir := filepath.Join(strings.Split(f.Function, ".")[:1]...)
				file := filepath.Base(f.File)
				path := fmt.Sprintf(" %s:%d", filepath.Join(dir, file), f.Line)
				return "", path
			},
		})
	case config.EnvProd:
		l.SetLevel(log.InfoLevel)
		l.SetFormatter(&log.JSONFormatter{})
	case config.EnvDev:
		l.SetLevel(log.DebugLevel)
		l.SetFormatter(&log.JSONFormatter{
			PrettyPrint: true,
		})
	default:
		panic(errb.Errorf("unknown env %s", env))
	}

	l.SetOutput(os.Stdout)
	l.SetReportCaller(true)
}

func mustSetupDB(cfg config.DB) *pgstor.DB {
	db, err := pgstor.Connect(cfg.DSN(), nil)
	if err != nil {
		log.Fatalln("failed to connect to postgres:", err)
	}
	return db
}

func mustRunServer(server *http.Server) {
	serveCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalln("failed to start server:", err)
		}
	}()

	<-serveCtx.Done()
	log.Infoln("shutting down server...")
	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalln("failed to shutdown server gracefully:", err)
	}
	log.Infoln("success")
}
