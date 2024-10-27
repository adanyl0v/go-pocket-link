package app

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"log"
)

var postgresConfig struct {
	Host     string `env:"POSTGRES_HOST"`
	Port     int    `env:"POSTGRES_PORT"`
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	DB       string `env:"POSTGRES_DATABASE"`
	SSLMode  string `env:"POSTGRES_SSL_MODE"`
}

func Run() {
	_ = godotenv.Load()
	err := cleanenv.ReadEnv(&postgresConfig)
	if err != nil {
		panic(err)
	}

	db := sqlx.MustConnect("pgx", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		postgresConfig.User, postgresConfig.Password, postgresConfig.Host, postgresConfig.Port,
		postgresConfig.DB, postgresConfig.SSLMode))
	defer func() { _ = db.Close() }()

	var users []struct {
		Email string `db:"email"`
	}
	err = db.Select(&users, "SELECT * FROM users")
	if err != nil {
		panic(err)
	}

	log.Println("Fetched", len(users), "users:")
	for _, user := range users {
		log.Printf("- %s;\n", user.Email)
	}
}
