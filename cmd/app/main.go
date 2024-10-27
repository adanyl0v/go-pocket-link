package main

import (
	_ "github.com/joho/godotenv/autoload"
	"go-pocket-link/internal/app"
	"os"
)

func main() {
	app.Run(os.Getenv("CONFIG_PATH"))
}
