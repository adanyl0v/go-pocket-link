package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"go-pocket-link/internal/config"
	"go-pocket-link/internal/domain"
	"go-pocket-link/internal/repository"
	"go-pocket-link/pkg/errb"
	pgstor "go-pocket-link/pkg/storage/postgres"
	"net/http"
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

	pgDB := mustSetupDB(cfg.DB)
	defer func() { _ = pgDB.Close() }()
	log.Infoln("connected to postgres")

	repos := repository.NewRepositories(pgDB)

	if cfg.Env != config.EnvLocal {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	if cfg.Env == config.EnvLocal {
		router.Use(gin.Logger())
	}

	api := router.Group("/api/v1")
	{
		repo := repos.Users
		usersGroup := api.Group("/users")
		usersGroup.GET("/", func(c *gin.Context) {
			users, err := repo.GetAll(c)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": fmt.Sprintf("failed to get users: %v", err),
				})
				return
			}

			c.JSON(http.StatusOK, users)
		})
		usersGroup.GET("/:id", func(c *gin.Context) {
			id, err := uuid.Parse(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": fmt.Sprintf("failed to get id: %v", err),
				})
				return
			}

			u := domain.User{ID: id}
			err = repo.GetByID(c, &u)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": fmt.Sprintf("failed to get user: %v", err),
				})
				return
			}

			c.JSON(http.StatusOK, u)
		})
		usersGroup.POST("/", func(c *gin.Context) {
			var u domain.User
			err := c.Bind(&u)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": fmt.Sprintf("failed to bind user: %v", err),
				})
				return
			}

			err = repo.Save(c, &u)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": fmt.Sprintf("failed to save user: %v", err),
				})
				return
			}

			c.Status(http.StatusCreated)
		})
		usersGroup.PUT("/:id", func(c *gin.Context) {
			id, err := uuid.Parse(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": fmt.Sprintf("failed to get id: %v", err),
				})
				return
			}

			var u domain.User
			err = c.Bind(&u)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": fmt.Sprintf("failed to bind user: %v", err),
				})
				return
			}

			u.ID = id
			err = repo.Update(c, &u)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": fmt.Sprintf("failed to update user: %v", err),
				})
				return
			}
		})
		usersGroup.DELETE("/:id", func(c *gin.Context) {
			id, err := uuid.Parse(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": fmt.Sprintf("failed to get id: %v", err),
				})
				return
			}

			err = repo.Delete(c, id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": fmt.Sprintf("failed to delete user: %v", err),
				})
				return
			}
		})
	}

	log.Infof("listening to %s:%d...\n", cfg.Server.Host, cfg.Server.Port)
	server := http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalln("failed to start server:", err)
	}
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
