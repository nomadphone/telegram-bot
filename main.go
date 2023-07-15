package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nomadphone/telegram-bot/handlers"
)

func runTelegramBot() {
	telegramToken := os.Getenv("TELEGRAM_TOKEN")
	bot, err := telegram.NewBotAPI(telegramToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	_, err = bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}

	updates := bot.ListenForWebhook("/" + bot.Token)
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	fmt.Println("Listening on port", port)
	go http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	handlers.HandleUpdates(bot, updates)
}

func main() {
	runTelegramBot()
}
