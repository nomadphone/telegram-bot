package main

import (
	"log"
	"os"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nomadphone/telegram-bot/users"
)

func runTelegramBot() {
	telegramToken := os.Getenv("TELEGRAM_TOKEN")
	bot, err := telegram.NewBotAPI(telegramToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := telegram.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("Got message from [%s] %s", update.Message.From.UserName, update.Message.Text)
			var reply string
			if update.Message.Text == "/register" {
				if users.IsInAllowList(update.Message.From.UserName) {
					user := users.ManageTelegramUser(update.Message.From, update.Message.Chat.ID)
					reply = "You have been registered. You can now send messages to your phone number and they will get redirected to your telegram user.\n"
					reply += "Your phone number is: " + user.TwillioPhone
				} else {
					reply = "You're not in the allow list. Please contact nomadphone@behindtherequests.com to get access."
				}
			} else {
				reply = "Unrecognized command. Please use /register to register your telegram user and start using your new number."
			}
			msg := telegram.NewMessage(update.Message.Chat.ID, reply)
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
	}
}

func main() {
	runTelegramBot()
}
