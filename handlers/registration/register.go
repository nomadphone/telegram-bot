package registration

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nomadphone/telegram-bot/users"
)

func Register(update telegram.Update) string {
	var reply string
	if users.IsInAllowList(update.Message.From.UserName) {
		user := users.ManageTelegramUser(update.Message.From, update.Message.Chat.ID)
		reply = "You have been registered. You can now send messages to your phone number and they will get redirected to your telegram user.\n"
		reply += "Your phone number is: " + user.TwillioPhone
	} else {
		reply = "You're not in the allow list. Please contact nomadphone@behindtherequests.com to get access."
	}
	return reply
}
