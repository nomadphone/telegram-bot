package handlers

import (
	"log"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nomadphone/telegram-bot/handlers/registration"
)

func HandleUpdates(bot *telegram.BotAPI, updates telegram.UpdatesChannel) {
	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("Got message from [%s] %s", update.Message.From.UserName, update.Message.Text)
			var reply string
			if update.Message.Text == "/register" {
				reply = registration.Register(update)
			} else {
				reply = "Unrecognized command. Please use /register to register your telegram user and start using your new number."
			}
			msg := telegram.NewMessage(update.Message.Chat.ID, reply)
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
	}
}
