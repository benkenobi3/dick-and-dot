package main

import (
	"context"
	"log"
	"os"

	botapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	dotenv "github.com/joho/godotenv"

	"github.com/benkenobi3/dick-and-dot/internal/telegram/update"
)

func main() {
	err := dotenv.Overload()
	if err != nil {
		log.Println(err)
	}

	ctx := context.Background()

	token, exists := os.LookupEnv("TELEGRAM_TOKEN")
	if !exists {
		log.Fatal("Cannot get telegram token from 'TELEGRAM_TOKEN' env variable")
	}

	databaseURL, exists := os.LookupEnv("DATABASE_URL")
	if !exists {
		log.Fatal("Cannot get database url from 'DATABASE_URL' env variable")
	}

	bot, err := botapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	db, err := sqlx.Connect("pgx", databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	updateConfig := botapi.NewUpdate(0)
	updateConfig.Timeout = 150

	updates := bot.GetUpdatesChan(updateConfig)

	updatesHandler := update.NewHandler(db, bot)
	for u := range updates {
		updatesHandler.Handle(ctx, u)
	}
}
