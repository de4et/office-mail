package main

import (
	"log"

	"github.com/de4et/office-mail/services/mail-worker/internal/app"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	app.Run(app.LoadConfigFromENV())
}
