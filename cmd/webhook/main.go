package main

import (
	"fmt"
	"log"
	"os"

	botapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	dotenv "github.com/joho/godotenv"

	"github.com/benkenobi3/dick-and-dot/internal/telegram/update"
)

func main() {
	err := dotenv.Overload()
	if err != nil {
		log.Panic(err)
	}

	token, exists := os.LookupEnv("TELEGRAM_TOKEN")
	if !exists {
		log.Panic("Cannot get telegram token from 'TELEGRAM_TOKEN' env variable")
	}

	host, exists := os.LookupEnv("HOST")
	if !exists {
		log.Panic("Cannot get hostname from 'HOST' env variable")
	}

	bot, err := botapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	wh, _ := botapi.NewWebhook(fmt.Sprintf("%s/%s", host, token))

	_, err = bot.Request(wh)
	if err != nil {
		log.Fatal(err)
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}

	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	updates := bot.ListenForWebhook("/" + token)

	updatesHandler := update.NewHandler(bot)
	for u := range updates {
		updatesHandler.Handle(u)
	}
}
