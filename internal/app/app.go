package app

import (
	"fmt"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"go-pocket-link/internal/config"
	"go-pocket-link/pkg/errb"
	pgstor "go-pocket-link/pkg/storage/postgres"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func Run(configPath string) {
	cfg := mustReadConfig(config.NewFileReader(configPath))
	mustSetupLogger(cfg.Env)

	log.Infof("read config '%s'\n", configPath)
	log.Infof("set up logger with '%s' env\n", cfg.Env)

	pgDB := mustSetupPostgres(cfg.Storage.Postgres)
	defer func() { _ = pgDB.Close() }()
	log.Infoln("connected to postgres")
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

func mustSetupPostgres(pgCfg config.Postgres) *pgstor.Storage {
	db, err := pgstor.New(pgCfg.DSN(), nil)
	if err != nil {
		log.Fatalln("failed to connect to postgres:", err)
	}
	return db
}
