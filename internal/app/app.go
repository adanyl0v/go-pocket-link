package app

import (
	"fmt"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"go-pocket-link/internal/config"
	"go-pocket-link/pkg/errb"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func Run(configPath string) {
	cfg := mustReadConfig(config.NewFileReader(configPath))
	mustSetupLogger(cfg.Env)

	log.Debugf("read config at '%s'", configPath)
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
