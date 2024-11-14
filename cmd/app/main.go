package main

import (
	"github.com/adanyl0v/go-pocket-link/internal/app"
	_ "github.com/joho/godotenv/autoload"
	"os"
)

func main() {
	app.Run(os.Getenv("CONFIG_PATH"))
}
