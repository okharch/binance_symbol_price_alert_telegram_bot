package main

import (
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func helpHandler(bot *tgbotapi.BotAPI, message *tgbotapi.Message, _ *sql.DB) error {
	helpText := "/alert symbol up/down price - Create a new price alert\n" +
		"/alerts - List all of your active alerts\n" +
		"/drop_all - Remove all of your active alerts\n" +
		"/help - Display this help message"

	msg := tgbotapi.NewMessage(message.Chat.ID, helpText)
	_, err := bot.Send(msg)
	if err != nil {
		return fmt.Errorf("error sending help message to Telegram: %s", err)
	}

	return nil
}
